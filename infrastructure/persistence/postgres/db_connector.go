package postgres

import (
	"fmt"
	"github.com/AeroAgency/go-admin-api/interfaces/rest/dto/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

// Объект для работы с БД
type DatabaseConnector struct {
	DB *gorm.DB
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
	db.Table(dto.ModelCode).Select("count(DISTINCT id)").Where(fmt.Sprintf("%s IS NOT NULL", dto.ModelFieldCode)).Limit(1).Count(&count)
	db = db.Table(dto.ModelCode).
		Select(fmt.Sprintf("DISTINCT(%s) as value", dto.ModelFieldCode)).
		Where(fmt.Sprintf("%s IS NOT NULL", dto.ModelFieldCode)).
		Limit(dto.Limit).Offset(dto.Offset).
		Scan(&rows.Items)
	rows.Total = count
	return &rows, db.Error
}

func (s *DatabaseConnector) GetModelFilterModelRefValues(dto models.ModelFilterValuesParamsDto) (*models.ValueRows, error) {
	var rows models.ValueRows
	table := dto.ModelFieldModelCode
	var count int
	// Подсчет total
	db := s.DB
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
