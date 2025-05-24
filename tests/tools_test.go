package tests

import (
	"update-service/internal/model"
	"update-service/internal/utils"
)

const (
	SuccessAllURL               = "http://127.0.0.1:40"
	SuccessUpdateOnlyMalwareURL = "http://127.0.0.1:50"
	SuccessUpdateOnlyRuleURL    = "http://127.0.0.1:60"

	// failed on checker worker
	FailedLoginURL       = "http://127.0.0.1:70"
	FailedStatusURL      = "http://127.0.0.1:80"
	FailedSoftVersionURL = "http://127.0.0.1:90"
	// failed on provide worker
	FailedProvideMalwareURL = "http://127.0.0.1:100"
	FailedProvideRulesURL   = "http://127.0.0.1:110"
	// failed on apply worker
	FailedApplyMalwareURL = "http://127.0.0.1:120"
	FailedApplyRulesURL   = "http://127.0.0.1:130"
	// grpc testing
	SuccessGRPCURL  = "http://127.0.0.1:140"
	SuccessGRPCUUID = "a5275643-b274-451d-848c-8cf6460384cb"
	FailedGRPCUUID  = "859e5742-d9cb-42c2-aad0-38d2edfb3396"
	FailedGRPCURL   = "http://127.0.0.1:150"

	//
	SuccessAllVersion           = "3.4"
	FailedProvideMalwareVersion = "12345"
	FailedProvideRulesVersion   = "123456"
)

var SuccessAllStatus = []model.Status{
	model.Status{
		Msg:    "",
		Status: "",
		Code:   utils.ERulesExpires,
	},
	model.Status{
		Msg:    "",
		Status: "",
		Code:   utils.EMalwareBaseExpires,
	},
}

var SuccessOnlyMalwareStatus = []model.Status{
	model.Status{
		Msg:    "",
		Status: "",
		Code:   utils.EMalwareBaseExpires,
	},
}

var SuccessOnlyRulesStatus = []model.Status{
	model.Status{
		Msg:    "",
		Status: "",
		Code:   utils.ERulesExpires,
	},
}
