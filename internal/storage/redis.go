package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func PingRedis(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}