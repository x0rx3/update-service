package tests

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"
	grpcLib "update-service/internal/grpc"
	"update-service/internal/grpc/gen"
	"update-service/internal/logging"
	"update-service/internal/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestMain(t *testing.T) {
	var (
		wg             sync.WaitGroup
		ctx, cancel    = context.WithCancel(context.Background())
		testDelay      = time.Duration(10) * time.Second
		requestTimeout = time.Duration(10) * time.Second
		sleepTimeDelay = time.Duration(2) * time.Second
	)
	log := logging.InitLogger("debug")
	updateChecker := service.NewUpdateChecker(log, NewIdsClientTest(), NewResultTableTest(), NewServerTableTest(), 3)
	updateProvider := service.NewUpdateProvider(log, NewUpdateServerClientTest(), NewResultTableTest(), NewServerTableTest(), "./cache", 3)
	updateApplier := service.NewUpdateApplier(log, NewIdsClientTest(), NewResultTableTest(), NewServerTableTest(), 3)

	service.NewWorkerManager(3, []service.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	service.NewPipeline(log, []service.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	producer := service.NewProduceManager(log, NewServerTableTest(), updateChecker.InputChan(), testDelay, 3)
	producer.Produce(&wg, ctx)

	grpcLib.NewServer(
		log,
		service.NewCheckServerGRPCService(
			service.NewCheckService(NewServerTableTest()),
			producer.InputChan(),
			requestTimeout,
		),
	).Start("127.0.0.1:9090", &wg, ctx)

	client, err := grpc.NewClient("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}

	updateCheckerClient := gen.NewUpdateCheckerClient(client)
	stream, err := updateCheckerClient.CheckUpdate(ctx, &gen.CheckUdateRequest{ServerUuid: uuid.NewString()})
	if err != nil {
		t.Error(err)
	}

	for {
		logMsg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Error(err)
		}

		if logMsg == nil {
			break
		}

		log.Debug("gRPC", zap.String("Title", logMsg.GetTitle()), zap.String("Description", logMsg.GetDescription()))
	}

	go func() {
		time.Sleep(sleepTimeDelay)
		cancel()
	}()

	wg.Wait()
	log.Info("Programm end work!")
}
