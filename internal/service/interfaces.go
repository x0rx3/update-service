package service

import (
	"context"
	"sync"
	"update-service/internal/model"
	"update-service/internal/utils"
)

type WorkerPull interface {
	Build(wg *sync.WaitGroup, ctx context.Context)
}

type Pipeline interface {
	Setup(wg *sync.WaitGroup, ctx context.Context)
}

type Producer interface {
	Produce(wg *sync.WaitGroup, ctx context.Context)
	InputChan() chan *model.Server
}

type Worker interface {
	Process(ctx context.Context)
	InputChan() chan *model.Task
	OutputChan() chan *model.Task
}

type IDSClient interface {
	Login(url, login, password string) error
	SoftVersion(url string) (string, error)
	Status(url string) ([]model.Status, error)
	Upload(idsUrl, filePath string, pkgType utils.PackageType) error
}

type UpdateServerClient interface {
	Login() error
	UpdateList(pkgType utils.PackageType) ([]model.RrUpdates, error)
	Download(pkgType utils.PackageType, pkgInfo *model.RrUpdates, dir4Save string) (string, error)
}

type Checker interface {
	Check(uuid string) (*model.Server, error)
}
