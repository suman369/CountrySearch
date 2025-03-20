package cache

import (
	"sync"
	"time"
)

// CacheItem represents an item in the cache with expiration
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// IsExpired checks if the cache item has expired
func (item CacheItem) IsExpired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Cache interface defines the operations for a cache
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
	Clear()
}

// InMemoryCache implements the Cache interface with thread-safe operations
type InMemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		items: make(map[string]CacheItem),
	}

	// Start a cleanup goroutine
	go cache.startCleanupTimer()

	return cache
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if item.IsExpired() {
		return nil, false
	}

	return item.Value, true
}

// Set adds a value to the cache
func (c *InMemoryCache) Set(key string, value interface{}) {
	c.SetWithExpiration(key, value, 0) // 0 means no expiration
}

// SetWithExpiration adds a value to the cache with expiration time
func (c *InMemoryCache) SetWithExpiration(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
}

// cleanup removes expired items from the cache
func (c *InMemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
}

// startCleanupTimer starts a timer to periodically clean up expired items
func (c *InMemoryCache) startCleanupTimer() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}
