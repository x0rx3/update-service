package tests

import (
	"update-service/internal/model"

	"github.com/google/uuid"
)

var SuccessAllServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessAllURL,
	Name:          "SuccessAllPKg-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessMalwareServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessUpdateOnlyMalwareURL,
	Name:          "SuccessMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessRulesServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessUpdateOnlyRuleURL,
	Name:          "SuccessMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedLoginServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedLoginURL,
	Name:          "FailedLogin-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedSoftVersionSever = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedSoftVersionURL,
	Name:          "FailedSoftVersion-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedProvideMalwareServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedProvideMalwareURL,
	Name:          "FailedProvideMalware",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedProvideRulesServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedProvideRulesURL,
	Name:          "FailedProvideRules-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedApplyMalwareServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedApplyMalwareURL,
	Name:          "FailedApplyMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedApplyRulesServer = &model.Server{
	UUID:          uuid.NewString(),
	Url:           FailedApplyRulesURL,
	Name:          "FailedApplyRUles-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessGRPCServer = &model.Server{
	UUID:          SuccessGRPCUUID,
	Url:           SuccessGRPCURL,
	Name:          "SuccessGRPC-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedGRPCServer = &model.Server{
	UUID:          FailedGRPCUUID,
	Url:           FailedGRPCURL,
	Name:          "SuccessGRPC-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
