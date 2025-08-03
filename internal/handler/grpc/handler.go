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

func NewHandler(accountHandler logic.AccountHandler) go_load.GoLoadServiceServer {
	return &Handler{
		accountHandler: accountHandler,
	}
}

func (h Handler) CreateAccount(ctx context.Context, request *go_load.CreateAccountRequest) (*go_load.CreateAccountResponse, error) {
	account, err := h.accountHandler.CreateAccount(ctx, logic.CreateAccountParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &go_load.CreateAccountResponse{
		AccountId: account.AccountID,
	}, nil
}

// CreateDownloadTask implements go_load.GoLoadServiceServer.
func (h *Handler) CreateDownloadTask(context.Context, *go_load.CreateDownloadTaskRequest) (*go_load.CreateDownloadTaskResponse, error) {
	panic("unimplemented")
}

// CreateSession implements go_load.GoLoadServiceServer.
func (h *Handler) CreateSession(context.Context, *go_load.CreateSessionRequest) (*go_load.CreateSessionResponse, error) {
	panic("unimplemented")
}

// DeleteDownloadTask implements go_load.GoLoadServiceServer.
func (h *Handler) DeleteDownloadTask(context.Context, *go_load.DeleteDownloadTaskRequest) (*go_load.DeleteDownloadTaskResponse, error) {
	panic("unimplemented")
}
func (h *Handler) GetDownloadTaskFile(*go_load.GetDownloadTaskFileRequest, grpc.ServerStreamingServer[go_load.GetDownloadTaskFileResponse]) error {
	panic("unimplemented")
}

// GetDownloadTaskList implements go_load.GoLoadServiceServer.
func (h *Handler) GetDownloadTaskList(context.Context, *go_load.GetDownloadTaskListRequest) (*go_load.GetDownloadTaskListResponse, error) {
	panic("unimplemented")
}

// UpdateDownloadTask implements go_load.GoLoadServiceServer.
func (h *Handler) UpdateDownloadTask(context.Context, *go_load.UpdateDownloadTaskRequest) (*go_load.UpdateDownloadTaskResponse, error) {
	panic("unimplemented")
}
