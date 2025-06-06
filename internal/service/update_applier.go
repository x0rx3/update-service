package service

import (
	"context"
	"fmt"
	"time"
	"update-service/internal/model"
	"update-service/internal/repository"
	"update-service/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UpdateApplier is responsible for applying updates to the server
type UpdateApplier struct {
	log             *zap.Logger            // Logger instance for debug and error logging
	vipNetIDSClient IDSClient              // Client to interact with IDS server
	resultTable     repository.ResultTable // Table to store update results
	serverTable     repository.ServerTable // Interface to store server in DB
	inputChan       chan *model.Task       // Channel to receive tasks for processing
	outputChan      chan *model.Task       // Channel to send tasks after processing
}

// NewUpdateApplier constructs a new instance of UpdateApplier
func NewUpdateApplier(
	log *zap.Logger,
	vipNetIDSClient IDSClient,
	resultTable repository.ResultTable,
	serverTable repository.ServerTable,
	limit int,
) *UpdateApplier {
	return &UpdateApplier{
		log:             log.With(zap.String("component", "UpdateApplier")),
		vipNetIDSClient: vipNetIDSClient,
		resultTable:     resultTable,
		serverTable:     serverTable,
		inputChan:       make(chan *model.Task, limit),
		outputChan:      make(chan *model.Task, limit),
	}
}

// Process listens to the input channel and starts processing tasks
func (inst *UpdateApplier) Process(ctx context.Context) {
	inst.log.Info("Start and Wait Task...")

	for {
		select {
		case job := <-inst.inputChan:
			job.SendProcessLog(&model.ProcessLog{Title: "Загрузка обновлений..."})
			inst.log.Info("Upload started", zap.String("server", job.Server().Name))

			result := inst.buildResult(job)

			inst.handle(job, result, utils.Malware)
			inst.handle(job, result, utils.Rules)

			job.SendProcessLog(&model.ProcessLog{Title: "Загрузка обновлений завершена"})
			inst.log.Info("Upload completed", zap.String("server", job.Server().Name))

			inst.complete(job, result)
		case <-ctx.Done():
			inst.log.Info("Shutdown signal received. Stopping...")
			return
		}
	}
}

// InputChan returns the input channel
func (inst *UpdateApplier) InputChan() chan *model.Task {
	return inst.inputChan
}

// OutputChan returns the output channel
func (inst *UpdateApplier) OutputChan() chan *model.Task {
	return inst.outputChan
}

// handle is responsible for managing the upload process of an update package (Malware or Rules)
// to a specified server. It logs each step of the process and updates the result and server status
// based on the outcome.
//
// The function performs the following steps:
// 1. Resolves the human-readable name of the package type.
// 2. Determines the file path for the update.
// 3. If the file path is invalid, logs the error and sets the result status.
// 4. If the file path is valid, attempts to upload the update file.
// 5. On success, updates the server and result statuses accordingly.
// 6. On failure, logs and reports the error.
func (inst *UpdateApplier) handle(job *model.Task, result *model.Result, pkgType utils.PackageType) {
	var pkgNameAlias string
	switch pkgType {
	case utils.Malware:
		pkgNameAlias = utils.MalwareNameAlias
	case utils.Rules:
		pkgNameAlias = utils.RulesNameAlias
	default:
		pkgNameAlias = "\"Unknown\""
	}

	filePath, err := inst.fileUpdatePath(job, pkgType)
	if err != nil {
		job.SendProcessLog(&model.ProcessLog{Title: fmt.Sprintf("Ошибка загрузки обновления %s", pkgNameAlias), Description: "Неверный формат пути до файла обновления"})
		inst.log.Error(fmt.Sprintf("Error process apply update %s", pkgNameAlias), zap.Error(err), zap.String("server", job.Server().Name))

		result.Errors += fmt.Sprintf("Failed upload %s: Invalid format path to update file", pkgNameAlias)

		if pkgType == utils.Malware {
			result.Malware = job.Server().MalwareStatus
		} else {
			result.Rules = job.Server().MalwareStatus
		}
	} else {
		if filePath != "" {
			inst.log.Info(fmt.Sprintf("Upload %s file", pkgNameAlias), zap.String("server", job.Server().Name))
			job.SendProcessLog(&model.ProcessLog{Title: fmt.Sprintf("Загрузка файла %s...", pkgNameAlias)})
			if inErr := inst.vipNetIDSClient.Upload(job.Server().Url, filePath, utils.Rules); inErr != nil {
				job.SendProcessLog(&model.ProcessLog{Title: fmt.Sprintf("Ошибка загрузки файла %s!", pkgNameAlias), Description: inErr.Error()})
				inst.log.Error(fmt.Sprintf("Error upload %s file", pkgNameAlias), zap.Error(err), zap.String("server", job.Server().Name))
			} else {
				job.SendProcessLog(&model.ProcessLog{Title: fmt.Sprintf("Загрузка файла %s успешно завершилась!", pkgNameAlias)})
				inst.log.Info(fmt.Sprintf("Upload %s file success completed", pkgNameAlias), zap.String("server", job.Server().Name))
				if pkgType == utils.Malware {
					job.Server().MalwareStatus = utils.UpdatedStatusSoftware
					result.Malware = utils.UpdatedStatusSoftware
				} else {
					job.Server().RulesStatus = utils.UpdatedStatusSoftware
					result.Rules = utils.UpdatedStatusSoftware
				}
			}
		} else {
			inst.log.Info(fmt.Sprintf("No update %s download required", pkgNameAlias), zap.String("server", job.Server().Name))
			job.SendProcessLog(&model.ProcessLog{Title: fmt.Sprintf("Файл обновления %s отсутствует", pkgNameAlias)})

		}
	}
}

func (inst *UpdateApplier) fileUpdatePath(job *model.Task, typ utils.PackageType) (string, error) {
	filePath, ok := job.Meta(string(typ))
	if !ok {
		return "", nil
	}

	filePathStr, ok := filePath.(string)
	if !ok {
		return "", fmt.Errorf("invalid path format ")
	}

	return filePathStr, nil
}

// buildResult initializes and returns a new Result for a job
func (inst *UpdateApplier) buildResult(job *model.Task) *model.Result {
	result := &model.Result{
		UUID:       uuid.NewString(),
		ServerUUID: job.Server().UUID,
		TimeEnd:    time.Now(),
		Malware:    job.Server().MalwareStatus,
		Rules:      job.Server().RulesStatus,
	}

	if timeStart, ok := job.Meta(utils.MetaTimeStart); ok {
		if val, isTime := timeStart.(time.Time); isTime {
			result.TimeStart = val
		}
	}
	return result
}

// complete finalizes the update process for a given job and result.
// It performs the following steps:
// 1. Sets the end time for the result.
// 2. Updates the server state in the repository.
// 3. Saves the result to the repository.
// 4. Sends logs to the job's process log channel to track progress and errors.
func (inst *UpdateApplier) complete(job *model.Task, result *model.Result) {
	result.TimeEnd = time.Now()
	job.SendProcessLog(&model.ProcessLog{Title: "Сохранения результата и обновление сервера БД"})
	if err := inst.serverTable.Update(job.Server()); err != nil {
		inst.log.Error("Error update server", zap.Error(err), zap.String("server", job.Server().Name))
		job.SendProcessLog(&model.ProcessLog{Title: "Ошибка обновления сервера в БД", Description: err.Error()})
	} else {
		job.SendProcessLog(&model.ProcessLog{Title: "Запись сервера обновлена"})
	}

	if resultUUID, err := inst.resultTable.Insert(result); err != nil {
		inst.log.Error("Error insert result", zap.Error(err), zap.String("server", job.Server().Name))
		job.SendProcessLog(&model.ProcessLog{Title: "Ошибка сохранения результата", Description: err.Error()})
	} else {
		job.SendProcessLog(&model.ProcessLog{Title: "Результат успешно сохранен", Description: fmt.Sprintf("Результат можно посмотреть по ID: %s", resultUUID)})
	}
	job.SendProcessLog(&model.ProcessLog{Title: "Процесс завершен"})
	job.SendProcessLog(nil)
}
