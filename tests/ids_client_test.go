package tests

import (
	"fmt"
	"update-service/internal/model"
	"update-service/internal/utils"
)

type IDSClientTest struct{}

func NewIdsClientTest() *IDSClientTest {
	return &IDSClientTest{}
}

func (inst *IDSClientTest) Login(url, login, password string) error {
	if url == FailedLoginURL {
		return fmt.Errorf("failed login")
	}
	return nil
}

func (inst *IDSClientTest) SoftVersion(url string) (string, error) {
	switch url {
	case FailedSoftVersionURL:
		return "", fmt.Errorf("failed get version")
	case SuccessAllURL:
		return SuccessAllVersion, nil
	case SuccessUpdateOnlyMalwareURL, SuccessUpdateOnlyRuleURL:
		return SuccessAllVersion, nil
	case FailedProvideMalwareURL:
		return FailedProvideMalwareVersion, nil
	case FailedProvideRulesURL:
		return FailedProvideRulesVersion, nil

	}
	return "3.9", nil
}

func (inst *IDSClientTest) Status(url string) ([]model.Status, error) {
	switch url {
	case FailedStatusURL:
		return []model.Status{}, fmt.Errorf("failed get status")
	case SuccessAllURL:
		return SuccessAllStatus, nil
	case SuccessUpdateOnlyMalwareURL:
		return SuccessOnlyMalwareStatus, nil
	case SuccessUpdateOnlyRuleURL:
		return SuccessOnlyRulesStatus, nil
	default:
		return SuccessAllStatus, nil
	}
}

func (inst *IDSClientTest) Upload(idsUrl, filePath string, pkgType utils.PackageType) error {
	if idsUrl == SuccessUpdateOnlyMalwareURL && pkgType == utils.Rules {
		return fmt.Errorf("failed upload file")
	}

	if idsUrl == SuccessUpdateOnlyRuleURL && pkgType == utils.Malware {
		return fmt.Errorf("failed upload file")
	}

	if idsUrl == FailedApplyMalwareURL && pkgType == utils.Malware {
		return fmt.Errorf("failed upload file")
	}

	if idsUrl == FailedApplyRulesURL && pkgType == utils.Rules {
		return fmt.Errorf("failed upload file")
	}

	return nil
}
