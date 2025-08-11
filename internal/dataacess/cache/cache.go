package cache

import (
	"context"
	"errors"
	"time"

	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, value ...any) error
	IsDataInSet(ctx context.Context, key string, value any) (bool, error)
}

type redisClient struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisClient(client *redis.Client, logger *zap.Logger) Cache {
	return &redisClient{client, logger}
}

func InitializeRedisClient(configs configs.CacheConfig, logger *zap.Logger) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.Address,
		Username: configs.UserName,
		Password: configs.Password,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		rdb.Close()
	}
	return rdb, cleanup, nil
}

func (r redisClient) Get(ctx context.Context, key string) (any, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("key does not exist")
		}
		return nil, err
	}
	return val, nil
}

func (r redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		r.logger.
			With(zap.String("key", key)).
			With(zap.Any("value", value)).
			With(zap.String("key", key), zap.Error(err)).
			Error("failed to set value in cache")
		return err
	}
	return nil
}

func (r redisClient) AddToSet(ctx context.Context, key string, value ...any) error {
	_, err := r.client.SAdd(ctx, key, value...).Result()
	if err != nil {
		r.logger.With(zap.String("key", key), zap.Any("value", value), zap.Error(err)).Error("failed to add value to set")
		return err
	}
	return nil
}

// IsDataInSet implements Cache.
func (r *redisClient) IsDataInSet(ctx context.Context, key string, value any) (bool, error) {
	result, err := r.client.SIsMember(ctx, key, value).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		r.logger.With(zap.String("key", key), zap.Any("value", value), zap.Error(err)).Error("failed to check if value is in set")
		return false, err
	}
	return result, nil
}
