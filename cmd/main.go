package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"
	"time"
	"update-service/internal/config"
	"update-service/internal/grpc"
	"update-service/internal/logging"
	"update-service/internal/service"
	"update-service/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "../config/config.yaml", "path to config.yaml file")
	flag.Parse()

	fmt.Println(uuid.NewString())
	fmt.Println(uuid.NewString())

	config, err := config.NewConfig(*configPath)
	if err != nil {
		fmt.Printf("failed read config: %s", err.Error())
		os.Exit(1)
	}

	log := logging.InitLogger(config.LogLevel)

	upodateDatabase, err := database.NewUpdateDatabase(log).Connect(config.DSN)
	if err != nil {
		log.Error("failed init connection db", zap.Error(err))
		os.Exit(1)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	idsClient := service.NewVIPNetIDSClient(client)
	updateServerClient := service.NewVIPNetUpdateServerClient(config.UpdateServerUrl, config.UpdateServerLogin, config.UpdateServerPassword, client)

	updateChecker := service.NewUpdateChecker(log, idsClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.WorkerLimit)
	updateProvider := service.NewUpdateProvider(log, updateServerClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.Cache, config.WorkerLimit)
	updateApplier := service.NewUpdateApplier(log, idsClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.WorkerLimit)

	var wg sync.WaitGroup
	ctx := context.Background()

	service.NewPipeline(log, []service.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	service.NewWorkerManager(config.WorkerLimit, []service.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)

	producer := service.NewProduceManager(log, upodateDatabase.ServerTable, updateChecker.InputChan(), time.Duration(config.Delay)*time.Hour, config.WorkerLimit)
	producer.Produce(&wg, ctx)

	grpc.NewServer(
		log,
		service.NewCheckServerGRPCService(
			service.NewCheckService(upodateDatabase.ServerTable),
			producer.InputChan(),
			time.Duration(3*time.Minute),
		),
	).Start(config.ServerAddress, &wg, ctx)

	wg.Wait()
}
