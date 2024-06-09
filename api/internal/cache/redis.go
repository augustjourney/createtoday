package cache

import (
	"context"
	"createtodayapi/internal/common"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {

	val, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return common.ErrCacheItemNotFound
	}

	if err != nil {
		return err
	}

	bs := []byte(val)

	err = json.Unmarshal(bs, dest)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Set(ctx context.Context, key string, val interface{}, exp *time.Duration) error {

	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}

	var expiration time.Duration = 0
	if exp != nil {
		expiration = *exp
	}

	err = r.client.Set(ctx, key, bs, expiration).Err()

	return err
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	return err
}

func (r *RedisCache) Reset(ctx context.Context) error {
	err := r.client.FlushAll(ctx).Err()
	return err
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}
