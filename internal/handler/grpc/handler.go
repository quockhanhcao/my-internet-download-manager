package grpc

import (
	"context"

	"github.com/quockhanhcao/my-internet-download-manager/internal/generated/grpc/go_load"
	"github.com/quockhanhcao/my-internet-download-manager/internal/logic"
	"google.golang.org/grpc"
)

type Handler struct {
	go_load.UnimplementedGoLoadServiceServer
	accountHandler logic.AccountHandler
}

func NewHandler() go_load.GoLoadServiceServer {
	return &Handler{}
}

func (h Handler) CreateAccount(ctx context.Context, request *go_load.CreateAccountRequest) (*go_load.CreateAccountResponse, error) {
	output, err := h.accountHandler.CreateAccount(
		ctx,
		logic.CreateAccountParams{
			AccountName: request.GetAccountName(),
			Password:    request.GetPassword(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &go_load.CreateAccountResponse{
		AccountId: output.ID,
	}, nil
}

func (h *Handler) CreateSession(context.Context, *go_load.CreateSessionRequest) (*go_load.CreateSessionResponse, error) {
	// Implement the logic for creating a session
	return nil, nil
}

func (h *Handler) CreateDownloadTask(context.Context, *go_load.CreateDownloadTaskRequest) (*go_load.CreateDownloadTaskResponse, error) {
	// Implement the logic for creating a download task
	return nil, nil
}

func (h *Handler) GetDownloadTaskList(context.Context, *go_load.GetDownloadTaskListRequest) (*go_load.GetDownloadTaskListResponse, error) {
	return nil, nil
}

func (h *Handler) UpdateDownloadTask(context.Context, *go_load.UpdateDownloadTaskRequest) (*go_load.UpdateDownloadTaskResponse, error) {
	return nil, nil
}

func (h *Handler) DeleteDownloadTask(context.Context, *go_load.DeleteDownloadTaskRequest) (*go_load.DeleteDownloadTaskResponse, error) {
	return nil, nil
}
func (h *Handler) GetDownloadTaskFile(*go_load.GetDownloadTaskFileRequest, grpc.ServerStreamingServer[go_load.GetDownloadTaskFileResponse]) error {
	return nil
}

func (h *Handler) mustEmbedUnimplementedGoLoadServiceServer() {
}
