package services

import (
	"context"
	"fmt"
	"os"
	"time"
	"update-service/pkg/database"
	"update-service/pkg/lib"
	"update-service/pkg/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpdateProvider handles the update process for malware and rules packages.
// It communicates with the remote server, checks if updates are needed,
// downloads the packages, and stores error results in the database.
type UpdateProvider struct {
	log                *zap.Logger          // Logger to record internal events and errors
	updateServerClient UpdateServerClient   // Client to interact with the update server
	resultTable        database.ResultTable // Interface to store processing results in DB
	serverTable        database.ServerTable // Interface to store server in DB
	inputChan          chan *models.Task    // Channel to receive incoming jobs
	outputChan         chan *models.Task    // Channel to forward jobs after processing
	cache              string               // Path to local cache directory
}

// NewUpdateProvider creates a new instance of UpdateProvider.
func NewUpdateProvider(
	log *zap.Logger,
	updateServerClient UpdateServerClient,
	resultTable database.ResultTable,
	serverTable database.ServerTable,
	cache string,
	limit int,
) *UpdateProvider {
	return &UpdateProvider{
		log:                log.With(zap.String("component", "UpdateProvider")),
		updateServerClient: updateServerClient,
		resultTable:        resultTable,
		serverTable:        serverTable,
		inputChan:          make(chan *models.Task, limit),
		outputChan:         make(chan *models.Task, limit),
		cache:              cache,
	}
}

// Process continuously listens for new jobs and handles them using provide().
func (inst *UpdateProvider) Process(ctx context.Context) {
	inst.log.Info("Start and Wait Task...")
	for {
		select {
		case job := <-inst.inputChan:
			job.SendProcessLog(&models.ProcessLog{Title: "Скачивание обновлений..."})
			inst.log.Info("Download started", zap.String("server", job.Server().Name))

			result := inst.buildResult(job)

			malwareErr := inst.handle(job, lib.Malware)
			rulesErr := inst.handle(job, lib.Rules)

			if malwareErr != nil && rulesErr != nil {
				job.SendProcessLog(&models.ProcessLog{Title: "Ошибка получения файлов обновления, процесс обновления прерван!"})
				inst.log.Info("Process donwload stoped, all files end with error", zap.String("server", job.Server().Name))
				inst.complete(job, result)
				continue
			}

			job.SendProcessLog(&models.ProcessLog{Title: "Скачивание обновлений завершено!"})
			inst.log.Info("Download completed", zap.String("server", job.Server().Name))
			inst.outputChan <- job

		case <-ctx.Done():
			inst.log.Info("Shutdown signal received. Stopping...")
			// close(inst.inputChan)
			// close(inst.outputChan)
			return
		}
	}
}

// InputChan exposes the input channel for external components to send jobs.
func (inst *UpdateProvider) InputChan() chan *models.Task {
	return inst.inputChan
}

// OutputChan exposes the output channel to retrieve processed jobs.
func (inst *UpdateProvider) OutputChan() chan *models.Task {
	return inst.outputChan
}

func (inst *UpdateProvider) handle(job *models.Task, pkgType lib.PackageType) error {
	var pkgNameAlias string
	switch pkgType {
	case lib.Malware:
		pkgNameAlias = lib.MalwareNameAlias
	case lib.Rules:
		pkgNameAlias = lib.RulesNameAlias
	default:
		pkgNameAlias = "\"Unknown\""
	}

	rulesFile, err := inst.requiredFile(job, lib.Rules)
	if rulesFile != "" && err != nil {
		job.SendProcessLog(&models.ProcessLog{Title: fmt.Sprintf("Ошибка скачивания файла %s", pkgNameAlias), Description: err.Error()})
		inst.log.Error(fmt.Sprintf("Error download file %s", pkgType.String()), zap.Error(err), zap.String("server", job.Server().Name))

		if val, ok := job.Meta("Error"); ok {
			job.AddMeta("Error", fmt.Sprintf("%v %s", val, err.Error()))
		} else {
			job.AddMeta("Error", err.Error)
		}

		job.AddMeta("Error", err.Error)
		job.DeleteMeta(string(lib.Rules))
		return err
	} else if rulesFile != "" && err == nil {
		job.SendProcessLog(&models.ProcessLog{Title: fmt.Sprintf("Скачивание файла %s успешно завершено!", pkgNameAlias)})
		inst.log.Info(fmt.Sprintf("Download %s success completed", pkgType.String()), zap.String("server", job.Server().Name))
		job.AddMeta(string(lib.Rules), rulesFile)
	}
	return nil
}

// requiremntPkg processes a job by checking and downloading required packages.
func (inst *UpdateProvider) requiredFile(job *models.Task, typ lib.PackageType) (string, error) {
	if _, ok := job.Meta(typ.String()); ok {
		filePath, err := inst.getFile(job, typ)
		if err != nil {
			return "", err
		}
		return filePath, nil
	}

	return "", nil
}

// buildResult initializes and returns a new Result for a job
func (inst *UpdateProvider) buildResult(job *models.Task) *models.Result {
	result := &models.Result{
		UUID:       uuid.NewString(),
		ServerUUID: job.Server().UUID,
		TimeEnd:    time.Now(),
		Malware:    job.Server().MalwareStatus,
		Rules:      job.Server().RulesStatus,
	}

	if timeStart, ok := job.Meta(lib.MetaTimeStart); ok {
		if val, isTime := timeStart.(time.Time); isTime {
			result.TimeStart = val
		}
	}
	return result
}

// getFile retrieves the specified package (malware/rules) for the job.
func (inst *UpdateProvider) getFile(job *models.Task, pkgType lib.PackageType) (string, error) {
	updateList, err := inst.updateServerClient.UpdateList(pkgType)
	if err != nil {
		return "", err
	}

	// Find the most recent package that matches server software version
	var pkgInfo *models.RrUpdates
	for _, u := range updateList {
		if u.Latest {
			for _, s := range u.Sw {
				if s == job.Server().SoftVersion {
					pkgInfo = &u
				}
			}
		}
	}

	if pkgInfo == nil {
		return "", fmt.Errorf("can't find pkg with support version: %s, package type: %s", job.Server().SoftVersion, pkgType.String())
	}

	// Check if file exists in local cache
	var filePath string
	filePath, err = inst.checkCache(pkgInfo)
	if err != nil {
		return "", err
	}

	// Download the package if it's not cached
	if filePath == "" {
		inst.log.Info("Login in Update Server")
		if err = inst.updateServerClient.Login(); err != nil {
			return "", err
		}

		inst.log.Info("Donwload package", zap.String("version", job.Server().SoftVersion), zap.String("type", pkgType.String()))
		filePath, err = inst.updateServerClient.Download(pkgType, pkgInfo, inst.cache)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

// checkCache verifies if a package file already exists in the local cache directory.
func (inst *UpdateProvider) checkCache(pkgInfo *models.RrUpdates) (string, error) {
	cacheFiles, err := os.ReadDir(inst.cache)
	if err != nil {
		return "", err
	}

	for _, file := range cacheFiles {
		if pkgInfo.Name == file.Name() {
			return fmt.Sprintf("%s/%s", inst.cache, file.Name()), nil
		}
	}

	// Return empty path if not found
	return "", nil
}

func (inst *UpdateProvider) complete(job *models.Task, result *models.Result) {
	result.TimeEnd = time.Now()
	job.SendProcessLog(&models.ProcessLog{Title: "Сохранения результата и обновление сервера БД"})
	if err := inst.serverTable.Update(job.Server()); err != nil {
		inst.log.Error("Error update server", zap.Error(err), zap.String("server", job.Server().Name))
		job.SendProcessLog(&models.ProcessLog{Title: "Ошибка обновления сервера в БД", Description: err.Error()})
	} else {
		job.SendProcessLog(&models.ProcessLog{Title: "Запись сервера обновлена"})
	}

	if resultUUID, err := inst.resultTable.Insert(result); err != nil {
		inst.log.Error("Error insert result", zap.Error(err), zap.String("server", job.Server().Name))
		job.SendProcessLog(&models.ProcessLog{Title: "Ошибка сохранения результата", Description: err.Error()})
	} else {
		job.SendProcessLog(&models.ProcessLog{Title: "Результат успешно сохранен", Description: fmt.Sprintf("Результат можно посмотреть по ID: %s", resultUUID)})
	}
	job.SendProcessLog(&models.ProcessLog{Title: "Процесс завершен"})
	job.SendProcessLog(nil)
}
