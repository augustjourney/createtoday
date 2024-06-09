package infra

import (
	"context"
	"createtodayapi/internal/logger"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis(host string, port string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Error(ctx, "Error connecting to Redis:", "err", err.Error())
		return nil, err
	}

	logger.Info(ctx, "connected to redis", "response", pong)

	return rdb, nil
}
