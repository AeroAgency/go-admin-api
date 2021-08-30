package service

import (
	"fmt"
	"github.com/AeroAgency/go-admin-api/infrastructure/persistence/postgres"
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

func (ms ModelsService) PingModelField(modelCode string, modelFieldCode string) error {
	_, err := strconv.Atoi(modelFieldCode)
	if err == nil {
		return fmt.Errorf("wrong code")
	}
	err = ms.db.PingModelField(modelCode, modelFieldCode)
	if err != nil {
		return err
	}
	return nil
}
