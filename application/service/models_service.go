package service

import (
	"fmt"
	appErrors "github.com/AeroAgency/go-admin-api/infrastructure/errors"
	"github.com/AeroAgency/go-admin-api/infrastructure/persistence/postgres"
	"github.com/AeroAgency/go-admin-api/interfaces/rest/dto/models"
	pkgErrors "github.com/pkg/errors"
	"strconv"
)

type ModelsService struct {
	db *postgres.DatabaseConnector
}

func NewModelsService(db *postgres.DatabaseConnector) *ModelsService {
	return &ModelsService{db: db}
}

func (ms ModelsService) PingModel(modelCode string) error {
	err := ms.db.PingModel(modelCode)
	if err != nil {
		return err
	}
	return nil
}

func (ms ModelsService) PingModelField(modelCode string, modelFieldCode string, modelFiledAddParams models.ModelFieldAddParamsApiDto) error {
	_, err := strconv.Atoi(modelFieldCode)
	if err == nil {
		return fmt.Errorf("wrong code")
	}
	if modelFiledAddParams.Type == "referenceLink" {
		if modelFiledAddParams.Multiple == true {
			ref := fmt.Sprintf("%s_reflink_%s", modelCode, modelFiledAddParams.ReferenceCode)
			err := ms.db.PingModel(ref)
			if err != nil {
				return err
			}
			return nil
		}
		if modelFieldCode != modelFiledAddParams.ReferenceCode+"_reflink" {
			return appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("value of param code is invalid"))}
		}
	} else if modelFiledAddParams.Type == "modelLink" {
		if modelFiledAddParams.Multiple == true {
			ref := fmt.Sprintf("%s_modellink_%s", modelCode, modelFiledAddParams.ModelCode)
			err := ms.db.PingModel(ref)
			if err != nil {
				return err
			}
			return nil
		}
		err = ms.db.PingModelField(modelCode, "id")
		if err != nil {
			return err
		}
		err = ms.db.PingModelField(modelCode, "name")
		if err != nil {
			return err
		}
		if modelFieldCode != modelFiledAddParams.ModelCode+"_modellink" {
			return appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("value of param code is invalid"))}
		}
	}
	err = ms.db.PingModelField(modelCode, modelFieldCode)
	if err != nil {
		return err
	}
	return nil
}

func (ms ModelsService) GetModelFilterValues(dto models.ModelFilterValuesParamsDto) (*models.ValueRows, error) {
	switch {
	case dto.ModelFieldType == "string":
		filterValues, err := ms.db.GetModelFilterStringValues(dto)
		if err != nil {
			return nil, err
		}
		return filterValues, nil
	case dto.ModelFieldType == "modelLink":
		filterValues, err := ms.db.GetModelFilterModelRefValues(dto)
		if err != nil {
			return nil, err
		}
		return filterValues, nil
	default:
		return nil, appErrors.BadRequestError{Err: pkgErrors.WithStack(fmt.Errorf("value of param type is invalid"))}
	}
}

func (ms ModelsService) GetModelElementsList(dto models.ModelElementsListParamsApiDto) (*models.ModelElements, error) {
	modelElements, err := ms.db.GetModelElementsList(dto)
	if err != nil {
		return nil, err
	}
	return modelElements, nil
}
