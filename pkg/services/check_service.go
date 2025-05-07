package services

import (
	"update-service/pkg/database"
	"update-service/pkg/models"
)

type CheckService struct {
	serverTable database.ServerTable
}

func NewCheckService(serverTable database.ServerTable) *CheckService {
	return &CheckService{
		serverTable: serverTable,
	}
}

func (inst *CheckService) Check(uuid string) (*models.Server, error) {
	return inst.serverTable.SelectOne(uuid)
}
