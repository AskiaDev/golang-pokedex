package pokecache

import (
	"fmt"
	"time"
)

// Add adds a new entry to the cache
func (c *Cache) Add(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if _, ok := c.cache[key]; ok {
		return fmt.Errorf("key already exists")
	}


	c.cache[key] = cacheEntry{
		data: value,
		createdAt: time.Now(),
	}
	return nil
}

// Get retrieves an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.data, true
}

// reapLoop is a loop that reap the cache every interval
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for range ticker.C {
		c.reap(interval)
	}
}

// reap reap the cache every interval
func (c *Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	for key, entry := range c.cache {
		if time.Since(entry.createdAt) > interval {
			delete(c.cache, key)
		}
	}
}
