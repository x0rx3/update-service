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
	"update-service/pkg/configs"
	"update-service/pkg/database"
	"update-service/pkg/grpc"
	"update-service/pkg/logging"
	"update-service/pkg/services"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "../config/config.yaml", "path to config.yaml file")
	flag.Parse()

	config, err := configs.NewConfig(*configPath)
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

	idsClient := services.NewVIPNetIDSClient(client)
	updateServerClient := services.NewVIPNetUpdateServerClient(config.UpdateServerUrl, config.UpdateServerLogin, config.UpdateServerPassword, client)

	updateChecker := services.NewUpdateChecker(log, idsClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.WorkerLimit)
	updateProvider := services.NewUpdateProvider(log, updateServerClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.Cache, config.WorkerLimit)
	updateApplier := services.NewUpdateApplier(log, idsClient, upodateDatabase.ResultTable, upodateDatabase.ServerTable, config.WorkerLimit)

	var wg sync.WaitGroup
	ctx := context.Background()

	services.NewPipeline(log, []services.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)
	services.NewWorkerManager(config.WorkerLimit, []services.Worker{updateChecker, updateProvider, updateApplier}).Build(&wg, ctx)

	producer := services.NewProduceManager(log, upodateDatabase.ServerTable, updateChecker.InputChan(), time.Duration(config.Delay)*time.Hour, config.WorkerLimit)
	producer.Produce(&wg, ctx)

	grpc.NewServer(
		log,
		services.NewCheckServerGRPCService(
			services.NewCheckService(upodateDatabase.ServerTable),
			producer.InputChan(),
			time.Duration(3*time.Minute),
		),
	).Start(config.ServerAddress, &wg, ctx)

	wg.Wait()
}
