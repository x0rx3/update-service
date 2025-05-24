package tests

import (
	"update-service/internal/model"
)

type ServerTableTest struct{}

func NewServerTableTest() *ServerTableTest {
	return &ServerTableTest{}
}

func (inst *ServerTableTest) SelectOne(uuid string) (*model.Server, error) {
	if uuid == SuccessGRPCUUID {
		return SuccessGRPCServer, nil
	}
	return FailedGRPCServer, nil
}

func (inst *ServerTableTest) SelectAll() ([]model.Server, error) {
	return []model.Server{
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

func (inst *ServerTableTest) Update(server *model.Server) error {
	return nil
}
