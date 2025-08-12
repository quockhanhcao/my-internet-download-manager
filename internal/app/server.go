package app

import (
	"context"
	"syscall"

	"github.com/quockhanhcao/my-internet-download-manager/internal/handler/grpc"
	"github.com/quockhanhcao/my-internet-download-manager/internal/handler/http"
	"github.com/quockhanhcao/my-internet-download-manager/internal/utils"
	"go.uber.org/zap"
)

type Server struct {
	grpcServer grpc.Server
	httpServer http.Server
	logger     *zap.Logger
}

func NewServer(grpcServer grpc.Server, httpServer http.Server, logger *zap.Logger) *Server {
	return &Server{
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s Server) Start() {
	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("gRPC server stopped")
	}()
	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("HTTP server stopped")
	}()
	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
}
