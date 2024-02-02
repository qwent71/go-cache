package cache

import (
	"sync"
	"time"
)

type CacheItem[T any] struct {
	Value    T
	ExpireAt time.Time
}

type Cache[T any] struct {
	items map[string]CacheItem[T]
	ttl   time.Duration
	m     sync.RWMutex
}

func NewCache[T any](ttl time.Duration) *Cache[T] {
	c := &Cache[T]{
		items: make(map[string]CacheItem[T]),
		ttl:   ttl,
	}

	go func() {
		for {
			time.Sleep(time.Second * 10) // Invalidate every 10 seconds
			c.Invalidate()
		}
	}()

	return c
}

// Set adds a new key-value pair to the cache
func (c *Cache[T]) Set(key string, value T) {
	c.SetWithLifeDuration(key, value, c.ttl)
}

func (c *Cache[T]) SetWithLifeDuration(key string, value T, lifeDuration time.Duration) {
	c.m.Lock()
	defer c.m.Unlock()

	item := make([]T, 1)
	copy(item, []T{value})

	c.items[key] = CacheItem[T]{
		Value:    item[0],
		ExpireAt: time.Now().Add(lifeDuration),
	}
}

// Get retrieves a value of a specific type from the cache
func (c *Cache[T]) Get(key string) (T, bool) {
	c.m.RLock()
	defer c.m.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return *new(T), false
	}
	if time.Now().After(item.ExpireAt) {
		delete(c.items, key)
		return *new(T), false
	}

	// Updated cache item expiration time
	c.items[key] = CacheItem[T]{
		Value:    item.Value,
		ExpireAt: time.Now().Add(c.ttl),
	}

	return item.Value, true
}

// invalide cache
func (c *Cache[T]) Invalidate() {
	c.m.Lock()
	defer c.m.Unlock()

	for k, v := range c.items {
		if time.Now().After(v.ExpireAt) {
			delete(c.items, k)
		}
	}
}

func (c *Cache[T]) Delete(key string) {
	c.m.Lock()
	defer c.m.Unlock()

	delete(c.items, key)
}
