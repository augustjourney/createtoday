package cache

import (
	"createtodayapi/internal/common"
	"encoding/json"
	"sync"
	"time"
)

type CacheItem struct {
	Data   []byte
	Expiry time.Time
}

type MemoryCache struct {
	store map[string]CacheItem
	mu    sync.RWMutex
}

func (m *MemoryCache) Get(key string, dest interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.store[key]

	if !ok || (!item.Expiry.IsZero() && time.Now().After(item.Expiry)) {
		delete(m.store, key)
		return common.ErrCacheItemNotFound
	}

	err := json.Unmarshal(item.Data, dest)
	if err != nil {
		return nil
	}

	return err
}

func (m *MemoryCache) Set(key string, val interface{}, exp *time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}

	item := CacheItem{
		Data: bs,
	}

	if exp != nil {
		item.Expiry = time.Now().Add(*exp)
	}

	m.store[key] = item

	return nil
}

func (m *MemoryCache) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.store, key)
	return nil
}

func (m *MemoryCache) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[string]CacheItem)
	return nil
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]CacheItem),
	}
}
