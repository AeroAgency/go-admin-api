package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
