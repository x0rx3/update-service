package database

import (
	"update-service/pkg/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MySQlResultTable struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewMysqlResultTable(log *zap.Logger, db *gorm.DB) *MySQlResultTable {
	return &MySQlResultTable{
		log: log,
		db:  db,
	}
}

func (inst *MySQlResultTable) Insert(result *models.Result) (string, error) {
	if result.UUID == "" {
		result.UUID = uuid.NewString()
	}

	queryResult := inst.db.Create(result)
	if queryResult.Error != nil {
		inst.log.Error("Error insert into result db", zap.Error(queryResult.Error))
		return "", queryResult.Error
	}

	return result.UUID, nil
}
