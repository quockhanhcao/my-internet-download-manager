package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/quockhanhcao/my-internet-download-manager/internal/generated/grpc/go_load"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
}

func NewServer() Server {
	return &server{}
}

func (s *server) Start(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := go_load.RegisterGoLoadServiceHandlerFromEndpoint(ctx, mux, "/api", opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8081", mux)
}
