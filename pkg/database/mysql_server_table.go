package database

import (
	"update-service/pkg/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MySqlServerTable struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewMySqlServerTable(log *zap.Logger, db *gorm.DB) *MySqlServerTable {
	return &MySqlServerTable{
		log: log,
		db:  db,
	}
}

func (inst *MySqlServerTable) SelectAll() ([]models.Server, error) {
	var servers []models.Server
	if result := inst.db.Find(&servers); result.Error != nil {
		return nil, result.Error
	}
	return servers, nil
}

func (inst *MySqlServerTable) SelectOne(uuid string) (*models.Server, error) {
	server := &models.Server{}
	if result := inst.db.First(server).Where("uuid = ?", uuid); result.Error != nil {
		return nil, result.Error
	}

	return server, nil
}

func (inst *MySqlServerTable) Update(server *models.Server) error {
	oldServer := &models.Server{}

	selectResult := inst.db.First(oldServer).Where("uuid = ?", server.UUID)
	if selectResult.Error != nil {
		return selectResult.Error
	}

	updatesField := map[string]any{}

	if oldServer.SoftVersion != server.SoftVersion {
		updatesField["soft_version"] = server.SoftVersion
	}

	if oldServer.MalwareStatus != server.MalwareStatus {
		updatesField["malware_status"] = server.MalwareStatus
	}

	if oldServer.RulesStatus != server.RulesStatus {
		updatesField["rules_status"] = server.RulesStatus
	}

	query := inst.db.Where("uuid = ?", server.UUID).Updates(updatesField)

	return query.Error
}
