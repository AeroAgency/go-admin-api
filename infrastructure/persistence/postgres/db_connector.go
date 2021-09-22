package postgres

import (
	"fmt"
	appErrors "github.com/AeroAgency/go-admin-api/infrastructure/errors"
	"github.com/AeroAgency/go-admin-api/interfaces/rest/dto/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	pkgErrors "github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// Объект для работы с БД
type DatabaseConnector struct {
	DB *gorm.DB
}

type LinkMultiFieldValues struct {
	Code   string
	Values []string
}

func NewDatabaseConnector(DB *gorm.DB) *DatabaseConnector {
	return &DatabaseConnector{DB: DB}
}

func (s *DatabaseConnector) PingModel(modelCode string) error {
	db := s.DB.Exec(fmt.Sprintf("SELECT * FROM %s LIMIT 1", modelCode))
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *DatabaseConnector) PingModelField(modelCode string, modelFieldCode string) error {
	db := s.DB.Exec(fmt.Sprintf("SELECT %s FROM %s LIMIT 1", modelFieldCode, modelCode))
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *DatabaseConnector) GetModelFilterStringValues(dto models.ModelFilterValuesParamsDto) (*models.ValueRows, error) {
	var rows models.ValueRows
	var count int
	// Подсчет total
	db := s.DB

	if len([]rune(dto.Query)) > 2 {
		db = db.Where(fmt.Sprintf("LOWER(%s) LIKE ?", dto.ModelFieldCode), "%"+strings.ToLower(dto.Query)+"%")
	}
	db.Table(dto.ModelCode).Select(fmt.Sprintf("count(DISTINCT %s)", dto.ModelFieldCode)).Where(fmt.Sprintf("%s IS NOT NULL", dto.ModelFieldCode)).Where(fmt.Sprintf("%s::varchar(255) !=''", dto.ModelFieldCode)).Limit(1).Count(&count)
	db = db.Table(dto.ModelCode).
		Select(fmt.Sprintf("DISTINCT(%s) as name", dto.ModelFieldCode)).
		Where(fmt.Sprintf("%s IS NOT NULL", dto.ModelFieldCode)).
		Where(fmt.Sprintf("%s::varchar(255) !=''", dto.ModelFieldCode)).
		Limit(dto.Limit).Offset(dto.Offset).
		Scan(&rows.Items)
	rows.Total = count
	return &rows, db.Error
}

func (s *DatabaseConnector) GetModelFilterModelRefValues(dto models.ModelFilterValuesParamsDto) (*models.ValueRows, error) {
	var rows models.ValueRows
	table := dto.ModelCode
	var count int
	// Подсчет total
	db := s.DB

	if len([]rune(dto.Query)) > 0 && len([]rune(dto.FilterId)) > 0 {
		return nil, appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("wrong input. only query or filter_id can be filled"))}
	}
	if len([]rune(dto.FilterId)) > 0 {
		ids := strings.Split(dto.FilterId, ",")
		db = db.Where("id IN (?)", ids)
	}
	if len([]rune(dto.Query)) > 2 {
		db = db.Where(fmt.Sprintf("%s LIKE ?", "LOWER(name)"), "%"+strings.ToLower(dto.Query)+"%")
	}
	db.Table(table).Select("count(DISTINCT id)").Where(fmt.Sprintf("%s IS NOT NULL", "name")).Limit(1).Count(&count)
	db = db.Table(table).
		Select("id as value, name").
		Where(fmt.Sprintf("%s IS NOT NULL", "name")).
		Limit(dto.Limit).Offset(dto.Offset).
		Scan(&rows.Items)
	rows.Total = count
	return &rows, db.Error
}

func (s *DatabaseConnector) GetModelElementsList(dto models.ModelElementsListParamsApiDto) (*models.ModelElements, error) {
	var rows models.ModelElements
	table := dto.ModelCode
	var count int
	// Подсчет total
	db := s.DB
	for _, v := range dto.Filter {
		field := v.Code
		values := v.Values
		isTypeMultiModelLink := strings.Contains(field, "_modellink_")
		isTypeMultiReferenceLink := strings.Contains(field, "_reflink_")
		if isTypeMultiModelLink {
			fieldData := strings.Split(field, "_modellink_")
			db = db.Where("id IN (SELECT DISTINCT ("+fieldData[0]+"_id) from "+v.Code+" WHERE "+fieldData[1]+"_id IN (?))", values)
			continue
		}
		if isTypeMultiReferenceLink {
			fieldData := strings.Split(field, "_reflink_")
			db = db.Where("id::text IN (SELECT DISTINCT ("+fieldData[0]+"_id) from "+v.Code+" WHERE "+fieldData[1]+"_id IN (?))", values)
			continue
		}
		if len(values) > 1 { // через OR
			db = db.Where(fmt.Sprintf("%s IN (?)", field), values)
		} else {
			value := values[0]
			db = db.Where(fmt.Sprintf("%s = '%s'", field, value))
		}
	}
	db.Table(table).Select("count(DISTINCT id)").Limit(1).Count(&count)
	db = db.Table(table).
		Select(dto.Select).
		Order(dto.Sort + " " + dto.Order + ", id desc").
		Limit(dto.Limit).Offset(dto.Offset).
		Scan(&rows.Items)
	rows.Total = count
	return &rows, db.Error
}

func (s DatabaseConnector) GetModelElement(modelCode string, modelElementId string, selectFields []string) (models.ModelElementDetail, error) {
	db := s.DB
	var simpleFields []string
	var linkedModelMultiFields []string
	var linkedRefMultiFields []string
	for _, field := range selectFields {
		isTypeMultiModelLink := strings.Contains(field, "_modellink_")
		isTypeMultiReferenceLink := strings.Contains(field, "_reflink_")
		if isTypeMultiModelLink {
			linkedModelMultiFields = append(linkedModelMultiFields, field)
		} else if isTypeMultiReferenceLink {
			linkedRefMultiFields = append(linkedRefMultiFields, field)
		} else {
			simpleFields = append(simpleFields, field)
		}
	}
	rows, err := db.Table(modelCode).
		Select(simpleFields).
		Where("id = ?", modelElementId).
		Limit(1).
		Rows()

	if err != nil {
		return nil, appErrors.NotFoundError{Err: pkgErrors.WithStack(err)}
	}
	cols, _ := rows.Columns()
	results := make([]map[string]interface{}, 1)
	for i := 0; rows.Next(); i++ {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}
		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		results[i] = m
		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
	}
	element := results[0]
	if element == nil {
		err = appErrors.NotFoundError{Err: pkgErrors.WithStack(fmt.Errorf("element not found"))}
		return nil, err
	}
	rowValues := make(map[string]string)
	for key, value := range element {
		if value == nil {
			rowValues[key] = ""
		} else {
			rowValues[key] = fmt.Sprintf("%v", value)
		}
	}
	if len(linkedModelMultiFields) > 0 {
		for _, linkedModelMultiField := range linkedModelMultiFields {
			fieldData := strings.Split(linkedModelMultiField, "_modellink_")
			var modelElementLinks []models.ModelElementLink
			db = db.Table(linkedModelMultiField).
				Select(fieldData[1]+"_id as link_id").
				Where(fieldData[0]+"_id = ?", modelElementId).
				Find(&modelElementLinks)
			for _, link := range modelElementLinks {
				val, ok := rowValues[linkedModelMultiField]
				if !ok {
					rowValues[linkedModelMultiField] = link.LinkId
				} else {
					fieldValues := []string{val}
					fieldValues = append(fieldValues, link.LinkId)
					rowValues[linkedModelMultiField] = strings.Join(fieldValues, ",")
				}
			}
		}
	}
	if len(linkedRefMultiFields) > 0 {
		for _, linkedRefMultiField := range linkedRefMultiFields {
			fieldData := strings.Split(linkedRefMultiField, "_reflink_")
			var modelElementLinks []models.ModelElementLink
			db = db.Table(linkedRefMultiField).
				Select(fieldData[1]+"_id as link_id").
				Where(fieldData[0]+"_id = ?", modelElementId).
				Find(&modelElementLinks)
			for _, link := range modelElementLinks {
				val, ok := rowValues[linkedRefMultiField]
				if !ok {
					rowValues[linkedRefMultiField] = link.LinkId
				} else {
					fieldValues := []string{val}
					fieldValues = append(fieldValues, link.LinkId)
					rowValues[linkedRefMultiField] = strings.Join(fieldValues, ",")
				}
			}
		}
	}
	return rowValues, nil
}

func (s DatabaseConnector) CreateModelElement(modelCode string, dto models.ModelElementCreateApiDto) (string, error) {
	db := s.DB
	id := uuid.NewV4().String()
	fields := []string{"id"}
	values := []string{id}

	var linkedModelMultiFields []LinkMultiFieldValues
	var linkedRefMultiFields []LinkMultiFieldValues
	for _, v := range dto.FieldValues {
		field := v.Code
		isTypeMultiModelLink := strings.Contains(field, "_modellink_")
		isTypeMultiReferenceLink := strings.Contains(field, "_reflink_")
		if isTypeMultiModelLink {
			linkedModelMultiFields = append(linkedModelMultiFields, LinkMultiFieldValues{Code: field, Values: v.Values})
		} else if isTypeMultiReferenceLink {
			linkedRefMultiFields = append(linkedRefMultiFields, LinkMultiFieldValues{Code: field, Values: v.Values})
		} else {
			fields = append(fields, field)
			values = append(values, strings.Join(v.Values, ","))
		}
	}
	db = db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (?)",
			modelCode,
			strings.Join(fields, ","),
		),
		values,
	)
	if db.Error != nil {
		if err, ok := db.Error.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return "", appErrors.ConflictError{Err: pkgErrors.WithStack(err)}
		}
		return "", db.Error
	}
	if len(linkedModelMultiFields) > 0 {
		for _, v := range linkedModelMultiFields {
			fieldData := strings.Split(v.Code, "_modellink_")
			fields := []string{
				fieldData[0] + "_id",
				fieldData[1] + "_id",
			}
			for _, vf := range v.Values {
				db = db.Exec(
					fmt.Sprintf(
						"INSERT INTO %s (%s) VALUES (?)",
						v.Code,
						strings.Join(fields, ","),
					),
					[]string{id, vf},
				)
				if db.Error != nil {
					return "", db.Error
				}
			}
		}
	}
	if len(linkedRefMultiFields) > 0 {
		for _, v := range linkedRefMultiFields {
			fieldData := strings.Split(v.Code, "_reflink_")
			fields := []string{
				fieldData[0] + "_id",
				fieldData[1] + "_id",
			}
			for _, vf := range v.Values {
				db = db.Exec(
					fmt.Sprintf(
						"INSERT INTO %s (%s) VALUES (?)",
						v.Code,
						strings.Join(fields, ","),
					),
					[]string{id, vf},
				)
				if db.Error != nil {
					return "", db.Error
				}
			}
		}
	}
	return id, nil
}

func (s DatabaseConnector) EditModelElement(modelCode string, modelElementId string, dto models.ModelElementCreateApiDto) (string, error) {
	db := s.DB
	var fields []string
	var values []string
	fieldValues := make(map[string]interface{})

	var linkedModelMultiFields []LinkMultiFieldValues
	var linkedRefMultiFields []LinkMultiFieldValues
	for _, v := range dto.FieldValues {
		field := v.Code
		isTypeMultiModelLink := strings.Contains(field, "_modellink_")
		isTypeMultiReferenceLink := strings.Contains(field, "_reflink_")
		if isTypeMultiModelLink {
			linkedModelMultiFields = append(linkedModelMultiFields, LinkMultiFieldValues{Code: field, Values: v.Values})
		} else if isTypeMultiReferenceLink {
			linkedRefMultiFields = append(linkedRefMultiFields, LinkMultiFieldValues{Code: field, Values: v.Values})
		} else {
			fields = append(fields, field)
			values = append(values, strings.Join(v.Values, ","))
			fieldValues[field] = strings.Join(v.Values, ",")
		}
	}
	db = db.Raw("WHERE id = ?", modelElementId)
	db = db.LogMode(true)
	db = db.Table(modelCode).Updates(fieldValues)
	if db.Error != nil {
		if err, ok := db.Error.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			return "", appErrors.ConflictError{Err: pkgErrors.WithStack(err)}
		}
		return "", db.Error
	}
	if len(linkedModelMultiFields) > 0 {
		for _, v := range linkedModelMultiFields {
			fieldData := strings.Split(v.Code, "_modellink_")
			fields := []string{
				fieldData[0] + "_id",
				fieldData[1] + "_id",
			}
			db = db.Exec(
				fmt.Sprintf(
					"DELETE FROM %s WHERE %s = '%s'",
					v.Code,
					fieldData[0]+"_id",
					modelElementId,
				),
			)
			for _, vf := range v.Values {
				db = db.Exec(
					fmt.Sprintf(
						"INSERT INTO %s (%s) VALUES (?)",
						v.Code,
						strings.Join(fields, ","),
					),
					[]string{modelElementId, vf},
				)
				if db.Error != nil {
					return "", db.Error
				}
			}
		}
	}
	if len(linkedRefMultiFields) > 0 {
		for _, v := range linkedRefMultiFields {
			fieldData := strings.Split(v.Code, "_reflink_")
			fields := []string{
				fieldData[0] + "_id",
				fieldData[1] + "_id",
			}
			db = db.Exec(
				fmt.Sprintf(
					"DELETE FROM %s WHERE %s = '%s'",
					v.Code,
					fieldData[0]+"_id",
					modelElementId,
				),
			)
			for _, vf := range v.Values {
				db = db.Exec(
					fmt.Sprintf(
						"INSERT INTO %s (%s) VALUES (?)",
						v.Code,
						strings.Join(fields, ","),
					),
					[]string{modelElementId, vf},
				)
				if db.Error != nil {
					return "", db.Error
				}
			}
		}
	}
	return modelElementId, nil
}

func (s DatabaseConnector) DeleteModelElement(modelCode string, modelElementId string) error {
	db := s.DB
	db = db.Table(modelCode).Where("id=?", modelElementId).Delete(modelElementId)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
