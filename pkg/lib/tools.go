package lib

import (
	"fmt"
	"strings"
)

const ERulesExpires = "ERulesExpires"
const EMalwareBaseExpires = "EMalwareBaseExpires"
const UpdatedStatusSoftware = "Updated"
const MetaCsrfToken = "X-Csrf-Token"
const MetaTimeStart = "Job-Time-Start"
const MalwareNameAlias = "\"Малвари\""
const RulesNameAlias = "\"Правил\""

const NotifyMessageDone = "Done"

type PackageType string

func (inst PackageType) String() string {
	return string(inst)
}

const Rules PackageType = "rr"
const Malware PackageType = "malware"

func EnsureTrailingSlash(url *string) {
	if strings.HasSuffix(*url, "/") {
		return
	}

	*url = fmt.Sprintf("%s/", *url)
}
