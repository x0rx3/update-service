package service

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"
	"update-service/internal/model"
	"update-service/internal/repository"
	"update-service/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	// Regex patterns to identify rules and malware-related status codes.
	regRules   = regexp.MustCompile(`(?i).*Rules.*`)
	regMalware = regexp.MustCompile(`(?i).*Malware.*`)
)

// UpdateChecker is responsible for checking whether an update is required on a server.
type UpdateChecker struct {
	log         *zap.Logger            // Logger to record internal events and errors
	resultTable repository.ResultTable // Interface to store processing results in DB
	serverTable repository.ServerTable // Interface to store server in DB
	idsClient   IDSClient              // Client to interact with the ids
	inputChan   chan *model.Task       // Channel to receive incoming jobs
	outputChan  chan *model.Task       // Channel to forward jobs after processing
	closeOnce   sync.Once
}

// NewUpdateChecker initializes and returns an UpdateChecker.
func NewUpdateChecker(log *zap.Logger, idsClient IDSClient, resTable repository.ResultTable, servTable repository.ServerTable, limit int) *UpdateChecker {
	return &UpdateChecker{
		log:         log.With(zap.String("component", "UpdateChecker")),
		resultTable: resTable,
		serverTable: servTable,
		idsClient:   idsClient,
		inputChan:   make(chan *model.Task, limit),
		outputChan:  make(chan *model.Task, limit),
	}
}

// Process continuously listens for incoming jobs and handles them.
func (inst *UpdateChecker) Process(ctx context.Context) {
	inst.log.Info("Start and Wait Task...")
	for {
		select {
		case job := <-inst.inputChan:
			job.SendProcessLog(&model.ProcessLog{Title: "Проверка необходимости обновления"})
			inst.log.Info("Checking update", zap.String("server", job.Server().Name))

			updated, err := inst.handleJob(job)
			if err != nil {
				job.SendProcessLog(&model.ProcessLog{Title: "Ошибка проверки необходимости обнволения", Description: err.Error()})
				inst.log.Error("Error process ckecking update", zap.Error(err), zap.String("server", job.Server().Name))
				inst.complete(job, err)
				continue
			}

			if updated {
				job.SendProcessLog(&model.ProcessLog{Title: "Обновление не требуется"})
				inst.log.Info("Already up to date", zap.String("server", job.Server().Name))
				inst.complete(job, nil)
				continue
			}

			inst.outputChan <- job
		case <-ctx.Done():
			inst.log.Info("Shutdown signal received. Stopping...")
			return
		}
	}
}

// InputChan returns the input channel for receiving jobs.
func (inst *UpdateChecker) InputChan() chan *model.Task {
	return inst.inputChan
}

// OutputChan returns the output channel for forwarding jobs.
func (inst *UpdateChecker) OutputChan() chan *model.Task {
	return inst.outputChan
}

// handleJob handles the complete update-checking logic for a single job.
func (inst *UpdateChecker) handleJob(job *model.Task) (bool, error) {
	job.AddMeta(utils.MetaTimeStart, time.Now())

	utils.EnsureTrailingSlash(&job.Server().Url)

	if err := inst.idsClient.Login(job.Server().Url, job.Server().Login, job.Server().Password); err != nil {
		inst.log.Error("Authorization error on the server", zap.Error(err))
		return false, fmt.Errorf("authorization error on the server")
	}

	status, err := inst.idsClient.Status(job.Server().Url)
	if err != nil {
		inst.log.Error("Failed to retrieve status", zap.Error(err))
		return false, fmt.Errorf("failed to retrieve status")
	}

	inst.evaluateStatus(status, job)

	// If update is required, proceed with version check.
	if inst.isUpdateRequired(job) {
		job.Server().SoftVersion, err = inst.idsClient.SoftVersion(job.Server().Url)
		if err != nil {
			inst.log.Error("Failed to retrieve status", zap.Error(err))
			return false, fmt.Errorf("failed to retrieve status")
		}
		return true, nil
	}

	return false, nil
}

func (inst *UpdateChecker) complete(job *model.Task, err error) {
	job.SendProcessLog(&model.ProcessLog{Title: "Сохранения результата и обновление сервера"})
	result := &model.Result{
		UUID:       uuid.NewString(),
		ServerUUID: job.Server().UUID,
		Malware:    job.Server().MalwareStatus,
		Rules:      job.Server().RulesStatus,
		TimeEnd:    time.Now(),
	}

	if err != nil {
		result.Errors = err.Error()
	}

	if t, ok := job.Meta(utils.MetaTimeStart); ok {
		result.TimeStart = t.(time.Time)
	}

	if err := inst.serverTable.Update(job.Server()); err != nil {
		inst.log.Error("Error update server", zap.Error(err), zap.String("server", job.Server().Name))
		job.SendProcessLog(&model.ProcessLog{Title: "Ошибка обновления сервера в БД", Description: err.Error()})
	} else {
		job.SendProcessLog(&model.ProcessLog{Title: "Запись в БД сервера обновлена"})
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

// evaluateStatus parses the status codes and updates job fields accordingly.
func (inst *UpdateChecker) evaluateStatus(status []model.Status, job *model.Task) {
	job.Server().RulesStatus = "updated"
	job.Server().MalwareStatus = "updated"

	for _, s := range status {
		switch s.Code {
		case utils.ERulesExpires:
			job.Server().RulesStatus = s.Code
			job.AddMeta(string(utils.Rules), nil)
		case utils.EMalwareBaseExpires:
			job.Server().MalwareStatus = s.Code
			job.AddMeta(string(utils.Malware), nil)
		default:
			if regRules.MatchString(s.Code) {
				job.Server().RulesStatus = s.Code
				continue
			}
			if regMalware.MatchString(s.Code) {
				job.Server().MalwareStatus = s.Code
				continue
			}
			inst.log.Warn("Unexpected status code", zap.String("server", job.Server().Name), zap.String("code", s.Code))
		}
	}
}

// isUpdateRequired checks whether the job requires an update based on server status.
func (inst *UpdateChecker) isUpdateRequired(job *model.Task) bool {
	return job.Server().RulesStatus == utils.ERulesExpires || job.Server().MalwareStatus == utils.EMalwareBaseExpires
}
