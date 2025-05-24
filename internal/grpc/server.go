package grpc

import (
	// "UpdateService/pkg/grpc"
	"context"
	"net"
	"sync"
	"update-service/internal/grpc/gen"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	log                 *zap.Logger
	updateCheckerServer gen.UpdateCheckerServer
}

func NewServer(log *zap.Logger, checkerService gen.UpdateCheckerServer) *Server {
	return &Server{
		log:                 log.With(zap.String("component", "gRPC Server")),
		updateCheckerServer: checkerService,
	}
}

func (inst *Server) Start(address string, wg *sync.WaitGroup, ctx context.Context) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()
	gen.RegisterUpdateCheckerServer(grpcServer, inst.updateCheckerServer)

	wg.Add(2)
	go func() {
		defer wg.Done()
		inst.log.Info("gRPC server start", zap.String("Address", address))
		if err := grpcServer.Serve(lis); err != nil {
			inst.log.Error("gRPC server failed. Restart service for new try", zap.Error(err))
			return
		}
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		inst.log.Info("Shutdown signal received. Stopping...")
		grpcServer.GracefulStop()
	}()

}
