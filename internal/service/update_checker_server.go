package service

import (
	"context"
	"time"
	"update-service/internal/grpc/gen"
	"update-service/internal/model"
)

type UpdateChekerServer struct {
	gen.UnimplementedUpdateCheckerServer
	produceChan   chan *model.Task
	service       Checker
	requesTimeout time.Duration
}

func NewCheckServerGRPCService(
	service Checker,
	produceChan chan *model.Task,
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

	processLogChan := make(chan *model.ProcessLog)
	defer close(processLogChan)
	ctx, cancel := context.WithTimeout(context.Background(), inst.requesTimeout)
	defer cancel()

	inst.produceChan <- model.NewTask(server, processLogChan)

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
