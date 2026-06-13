package memory

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	value []byte
	exp   time.Time
}

func NewCache() *Cache { return &Cache{data: map[string]entry{}} }

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[key]
	if !ok || (!e.exp.IsZero() && time.Now().After(e.exp)) {
		return nil, nil
	}
	return e.value, nil
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.data[key] = entry{value: append([]byte(nil), value...), exp: exp}
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
