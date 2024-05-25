package cache

import "time"

type Cache interface {
	Get(key string, dest interface{}) error
	Set(key string, val interface{}, exp *time.Duration) error
	Delete(key string) error
	Reset() error
}
