package mysql

import (
	"update-service/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewServer(log *zap.Logger, db *gorm.DB) *Server {
	return &Server{
		log: log,
		db:  db,
	}
}

func (inst *Server) SelectAll() ([]model.Server, error) {
	var servers []model.Server
	if result := inst.db.Find(&servers); result.Error != nil {
		return nil, result.Error
	}
	return servers, nil
}

func (inst *Server) SelectOne(uuid string) (*model.Server, error) {
	server := &model.Server{}
	if result := inst.db.First(server).Where("uuid = ?", uuid); result.Error != nil {
		return nil, result.Error
	}

	return server, nil
}

func (inst *Server) Update(server *model.Server) error {
	oldServer := &model.Server{}

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
