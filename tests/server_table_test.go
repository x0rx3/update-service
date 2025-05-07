package tests

import (
	"update-service/pkg/models"
)

type ServerTableTest struct{}

func NewServerTableTest() *ServerTableTest {
	return &ServerTableTest{}
}

func (inst *ServerTableTest) SelectOne(uuid string) (*models.Server, error) {
	if uuid == SuccessGRPCUUID {
		return SuccessGRPCServer, nil
	}
	return FailedGRPCServer, nil
}

func (inst *ServerTableTest) SelectAll() ([]models.Server, error) {
	return []models.Server{
		*SuccessAllServer,
		*SuccessMalwareServer,
		*SuccessRulesServer,
		*FailedLoginServer,
		*FailedSoftVersionSever,
		*FailedProvideMalwareServer,
		*FailedProvideRulesServer,
		*FailedApplyRulesServer,
		*FailedApplyMalwareServer,
	}, nil
}

func (inst *ServerTableTest) Update(server *models.Server) error {
	return nil
}
