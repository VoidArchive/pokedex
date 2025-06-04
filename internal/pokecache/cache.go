package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries      map[string]cacheEntry
	mu           sync.Mutex
	reapInterval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:      make(map[string]cacheEntry),
		reapInterval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		time.Now().UTC(),
		val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.reapInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	expirationTime := time.Now().UTC().Add(-c.reapInterval)

	for key, entry := range c.entries {
		if entry.createdAt.Before(expirationTime) {
			delete(c.entries, key)
		}
	}
}
