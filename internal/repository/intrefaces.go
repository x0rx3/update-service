package repository

import "update-service/internal/model"

type Database interface {
	Connect(dsn string) error
}

type ResultTable interface {
	Insert(result *model.Result) (string, error)
}

type ServerTable interface {
	SelectOne(uuid string) (*model.Server, error)
	SelectAll() ([]model.Server, error)
	Update(server *model.Server) error
}
