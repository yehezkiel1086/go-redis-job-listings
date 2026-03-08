package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/config"
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context, conf *config.Redis) (*Redis, error) {
	db, err := strconv.Atoi(conf.DB)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       db,
		Protocol: 2,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := r.client.Del(ctx, key).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
