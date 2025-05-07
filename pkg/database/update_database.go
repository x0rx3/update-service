package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UpdateDatabase struct {
	db          *gorm.DB
	log         *zap.Logger
	ResultTable ResultTable
	ServerTable ServerTable
}

func NewUpdateDatabase(log *zap.Logger) *UpdateDatabase {
	return &UpdateDatabase{
		log: log,
	}
}

func (inst *UpdateDatabase) Connect(dsn string) (*UpdateDatabase, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	inst.log.Info("Success conected", zap.String("component", "UpdateDatabase"))
	inst.db = db
	inst.ResultTable = NewMysqlResultTable(inst.log, inst.db)
	inst.ServerTable = NewMySqlServerTable(inst.log, inst.db)

	return inst, nil
}
