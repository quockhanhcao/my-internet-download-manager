package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type DownloadTask struct {
	OfAccountId    uint64 `sql:"of_account_id"`
	TaskID         uint64 `sql:"task_id"`
	DownloadType   string `sql:"download_type"`
	URL            string `sql:"url"`
	DownloadStatus int    `sql:"download_status"`
	Metadata       string `sql:"metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, task DownloadTask) error
	GetDownloadTaskByID(ctx context.Context, id uint64) (DownloadTask, error)
	GetDownloadTasksByAccountID(ctx context.Context, accountID uint64) ([]DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task DownloadTask) error
	DeleteDownloadTask(ctx context.Context, id uint64) error
}

type downloadTaskAccessor struct {
	database *goqu.Database
}

// CreateDownloadTask implements DownloadTaskDataAccessor.
func (a *downloadTaskAccessor) CreateDownloadTask(ctx context.Context, task DownloadTask) error {
	panic("unimplemented")
}

// GetDownloadTaskByID implements DownloadTaskDataAccessor.
func (a *downloadTaskAccessor) GetDownloadTaskByID(ctx context.Context, id uint64) (DownloadTask, error) {
	panic("unimplemented")
}

// GetDownloadTasksByAccountID implements DownloadTaskDataAccessor.
func (a *downloadTaskAccessor) GetDownloadTasksByAccountID(ctx context.Context, accountID uint64) ([]DownloadTask, error) {
	panic("unimplemented")
}

// UpdateDownloadTask implements DownloadTaskDataAccessor.
func (a *downloadTaskAccessor) UpdateDownloadTask(ctx context.Context, task DownloadTask) error {
	panic("unimplemented")
}

// DeleteDownloadTask implements DownloadTaskDataAccessor.
func (a *downloadTaskAccessor) DeleteDownloadTask(ctx context.Context, id uint64) error {
	panic("unimplemented")
}

func NewDownloadTaskDataAccessor(database *goqu.Database) DownloadTaskDataAccessor {
	return &downloadTaskAccessor{database}
}
