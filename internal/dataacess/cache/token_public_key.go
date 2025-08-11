package cache

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type TokenPublicKeyCache interface {
	GetTokenPublicKey(ctx context.Context, id uint64) ([]byte, error)
	SetTokenPublicKey(ctx context.Context, id uint64, value []byte) error
}

type tokenPublicKeyCache struct {
	client Cache
	logger *zap.Logger
}

func NewTokenPublicKeyCache(client Cache, logger *zap.Logger) TokenPublicKeyCache {
	return &tokenPublicKeyCache{
		client: client,
		logger: logger,
	}
}

func getCacheKey(id uint64) string {
	return fmt.Sprintf("token_public_key:%d", id)
}

// GetTokenPublicKey implements TokenPublicKeyCache.
func (t tokenPublicKeyCache) GetTokenPublicKey(ctx context.Context, id uint64) ([]byte, error) {
	cacheKey := getCacheKey(id)
	result, err := t.client.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	publicKey, ok := result.([]byte)
	if !ok {
		return nil, errors.New("cache entry is not of type bytes")
	}
	return publicKey, nil
}

// SetTokenPublicKey implements TokenPublicKeyCache.
func (t tokenPublicKeyCache) SetTokenPublicKey(ctx context.Context, id uint64, value []byte) error {
	cacheKey := getCacheKey(id)
	return t.client.Set(ctx, cacheKey, value, 0)
}
