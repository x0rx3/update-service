package services

import (
	"context"
	"time"
	"update-service/pkg/grpc/gen"
	"update-service/pkg/models"
)

type UpdateChekerServer struct {
	gen.UnimplementedUpdateCheckerServer
	produceChan   chan *models.Task
	service       Checker
	requesTimeout time.Duration
}

func NewCheckServerGRPCService(
	service Checker,
	produceChan chan *models.Task,
	requestTimeout time.Duration,
) *UpdateChekerServer {
	return &UpdateChekerServer{
		produceChan:   produceChan,
		service:       service,
		requesTimeout: requestTimeout,
	}
}

func (inst *UpdateChekerServer) CheckUpdate(req *gen.CheckUdateRequest, stream gen.UpdateChecker_CheckUpdateServer) error {
	server, err := inst.service.Check(req.ServerUuid)
	if err != nil {
		return err
	}

	processLogChan := make(chan *models.ProcessLog)
	defer close(processLogChan)
	ctx, cancel := context.WithTimeout(context.Background(), inst.requesTimeout)
	defer cancel()

	inst.produceChan <- models.NewTask(server, processLogChan)

	for {
		select {
		case log := <-processLogChan:
			if log == nil {
				return nil
			}
			respMsg := &gen.CheckUdateResponce{
				Title:       log.Title,
				Description: log.Description,
			}
			if err := stream.Send(respMsg); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
