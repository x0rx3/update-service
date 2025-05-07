package tests

import (
	"update-service/pkg/models"

	"github.com/google/uuid"
)

var SuccessAllServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessAllURL,
	Name:          "SuccessAllPKg-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessMalwareServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessUpdateOnlyMalwareURL,
	Name:          "SuccessMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessRulesServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           SuccessUpdateOnlyRuleURL,
	Name:          "SuccessMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedLoginServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedLoginURL,
	Name:          "FailedLogin-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedSoftVersionSever = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedSoftVersionURL,
	Name:          "FailedSoftVersion-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedProvideMalwareServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedProvideMalwareURL,
	Name:          "FailedProvideMalware",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedProvideRulesServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedProvideRulesURL,
	Name:          "FailedProvideRules-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedApplyMalwareServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedApplyMalwareURL,
	Name:          "FailedApplyMalware-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedApplyRulesServer = &models.Server{
	UUID:          uuid.NewString(),
	Url:           FailedApplyRulesURL,
	Name:          "FailedApplyRUles-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var SuccessGRPCServer = &models.Server{
	UUID:          SuccessGRPCUUID,
	Url:           SuccessGRPCURL,
	Name:          "SuccessGRPC-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
var FailedGRPCServer = &models.Server{
	UUID:          FailedGRPCUUID,
	Url:           FailedGRPCURL,
	Name:          "SuccessGRPC-Server",
	Login:         "login",
	Password:      "password",
	SoftVersion:   "",
	RulesStatus:   "",
	MalwareStatus: "",
}
