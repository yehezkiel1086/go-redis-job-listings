package port

import (
	"context"
	"time"
)

//go:generate mockgen -source=redis.go -destination=../mock/redis_mock.go -package=mock

type CacheRepository interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, keys ...string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
}
