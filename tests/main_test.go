package tests

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"
	grpcLib "update-service/pkg/grpc"
	"update-service/pkg/grpc/gen"
	"update-service/pkg/logging"
	"update-service/pkg/services"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestMain(t *testing.T) {
	var (
		wg             sync.WaitGroup
		ctx, cancel    = context.WithCancel(context.Background())
		testDelay      = time.Duration(10) * time.Minute
		requestTimeout = time.Duration(10) * time.Minute
		sleepTimeDelay = time.Duration(15) * time.Minute
	)
	log := logging.InitLogger("debug")
	updateChecker := services.NewUpdateChecker(log, NewIdsClientTest(), NewResultTableTest(), NewServerTableTest(), 3)
	updateProvider := services.NewUpdateProvider(log, NewUpdateServerClientTest(), NewResultTableTest(), NewServerTableTest(), "./cache", 3)
	updateApplier := services.NewUpdateApplier(log, NewIdsClientTest(), NewResultTableTest(), NewServerTableTest(), 3)

	services.NewWorkerManager(3, []services.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	services.NewPipeline(log, []services.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	producer := services.NewProduceManager(log, NewServerTableTest(), updateChecker.InputChan(), testDelay, 3)
	producer.Produce(&wg, ctx)

	grpcLib.NewServer(
		log,
		services.NewCheckServerGRPCService(
			services.NewCheckService(NewServerTableTest()),
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
