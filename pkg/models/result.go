package models

import "time"

type Result struct {
	UUID       string    `json:"uuid" gorm:"<-:create;column:uuid;primarykey"`
	ServerUUID string    `json:"server_uuid" gorm:"<-:create;column:server_uuid"`
	TimeStart  time.Time `json:"time_start" gorm:"<-:create;column:time_start"`
	TimeEnd    time.Time `json:"time_end" gorm:"<-:create;column:time_end"`
	Errors     string    `json:"errors" gorm:"<-:create;column:errors"`
	Malware    string    `json:"malware" gorm:"<-:create;column:malware"`
	Rules      string    `json:"rules" gorm:"<-:create;column:rules"`
}

func (inst Result) TableName() string {
	return "result"
}
