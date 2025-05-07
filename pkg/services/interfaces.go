package services

import (
	"context"
	"sync"
	"update-service/pkg/lib"
	"update-service/pkg/models"
)

type WorkerPull interface {
	Build(wg *sync.WaitGroup, ctx context.Context)
}

type Pipeline interface {
	Setup(wg *sync.WaitGroup, ctx context.Context)
}

type Producer interface {
	Produce(wg *sync.WaitGroup, ctx context.Context)
	InputChan() chan *models.Server
}

type Worker interface {
	Process(ctx context.Context)
	InputChan() chan *models.Task
	OutputChan() chan *models.Task
}

type IDSClient interface {
	Login(url, login, password string) error
	SoftVersion(url string) (string, error)
	Status(url string) ([]models.Status, error)
	Upload(idsUrl, filePath string, pkgType lib.PackageType) error
}

type UpdateServerClient interface {
	Login() error
	UpdateList(pkgType lib.PackageType) ([]models.RrUpdates, error)
	Download(pkgType lib.PackageType, pkgInfo *models.RrUpdates, dir4Save string) (string, error)
}

type Checker interface {
	Check(uuid string) (*models.Server, error)
}
