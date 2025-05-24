package database

import (
	"update-service/internal/repository"
	"update-service/internal/repository/mysql"

	"go.uber.org/zap"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UpdateDatabase struct {
	db          *gorm.DB
	log         *zap.Logger
	ResultTable repository.ResultTable
	ServerTable repository.ServerTable
}

func NewUpdateDatabase(log *zap.Logger) *UpdateDatabase {
	return &UpdateDatabase{
		log: log,
	}
}

func (inst *UpdateDatabase) Connect(dsn string) (*UpdateDatabase, error) {
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	inst.log.Info("Success conected", zap.String("component", "UpdateDatabase"))
	inst.db = db
	inst.ResultTable = mysql.NewResult(inst.log, inst.db)
	inst.ServerTable = mysql.NewServer(inst.log, inst.db)

	return inst, nil
}
