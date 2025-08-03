package grpc

import (
	"context"
	"net"

	"github.com/quockhanhcao/my-internet-download-manager/internal/generated/grpc/go_load"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler go_load.GoLoadServiceServer
}

func NewServer(handler go_load.GoLoadServiceServer) Server {
	return &server{handler: handler}
}

func (s *server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}

	defer listener.Close()

	grpcServer := grpc.NewServer()
	go_load.RegisterGoLoadServiceServer(grpcServer, s.handler)
	return grpcServer.Serve(listener)
}
