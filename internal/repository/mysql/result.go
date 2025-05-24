package mysql

import (
	"update-service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Result struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewResult(log *zap.Logger, db *gorm.DB) *Result {
	return &Result{
		log: log,
		db:  db,
	}
}

func (inst *Result) Insert(result *model.Result) (string, error) {
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
