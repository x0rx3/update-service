package models

type Server struct {
	UUID          string `json:"uuid" gorm:"<-:false;primaryKey;column:uuid"` // allow read, disable write permissions
	Url           string `json:"url" gorm:"<-:update;column:url"`
	Name          string `json:"name" gorm:"<-:false;column:name"`
	Login         string `json:"login" gorm:"<-:false;column:login"`
	Password      string `json:"password" gorm:"<-:false;column:password"`
	SoftVersion   string `json:"soft_version" gorm:"<-:update;column:soft_version"`
	RulesStatus   string `json:"rules_status" gorm:"<-:update;column:rules_status"`
	MalwareStatus string `json:"malware_status" gorm:"<-:update;column:malware_status"`
}

func (isnt Server) TableName() string {
	return "servers"
}
