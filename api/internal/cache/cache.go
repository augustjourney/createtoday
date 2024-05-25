package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, val interface{}, exp *time.Duration) error
	Delete(ctx context.Context, key string) error
	Reset(ctx context.Context) error
}
