package service

import (
	"update-service/internal/model"
	"update-service/internal/repository"
)

type Check struct {
	serverTable repository.ServerTable
}

func NewCheckService(serverTable repository.ServerTable) *Check {
	return &Check{
		serverTable: serverTable,
	}
}

func (inst *Check) Check(uuid string) (*model.Server, error) {
	return inst.serverTable.SelectOne(uuid)
}
