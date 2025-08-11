package cache

import (
	"context"

	"go.uber.org/zap"
)

const (
	AccountNameTakenSet = "account_name_taken"
)

type AccountNameCache interface {
	IsAccountNameTaken(ctx context.Context, name string) (bool, error)
	SetAccountNameTaken(ctx context.Context, name string) error
}

type accountNameCache struct {
	client Cache
	logger *zap.Logger
}

func (a accountNameCache) IsAccountNameTaken(ctx context.Context, name string) (bool, error) {
	return a.client.IsDataInSet(ctx, AccountNameTakenSet, name)
}

// SetAccountNameTaken implements AccountNameTakenCache.
func (a accountNameCache) SetAccountNameTaken(ctx context.Context, name string) error {
	return a.client.AddToSet(ctx, AccountNameTakenSet, name)
}

func NewAccountNameCache(client Cache, logger *zap.Logger) AccountNameCache {
	return &accountNameCache{
		client: client,
		logger: logger,
	}
}
