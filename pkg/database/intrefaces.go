package database

import "update-service/pkg/models"

type Database interface {
	Connect(dsn string) error
}

type ResultTable interface {
	Insert(result *models.Result) (string, error)
}

type ServerTable interface {
	SelectOne(uuid string) (*models.Server, error)
	SelectAll() ([]models.Server, error)
	Update(server *models.Server) error
}
