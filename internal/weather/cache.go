package weather

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value     string
	CreatedAt time.Time
}

type Cache struct {
	items map[string]CacheItem
	ttl   time.Duration
	mu    sync.RWMutex
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return "", false
	}
	if time.Since(item.CreatedAt) > c.ttl {
		return "", false
	}
	return item.Value, true
}

func (c *Cache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{Value: value, CreatedAt: time.Now()}
}
