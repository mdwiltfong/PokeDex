package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	time time.Time
	val  []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		time: time.Now(),
		val:  val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, exists := c.cache[key]; exists {
		return value.val, exists
	}
	return nil, false
}

func (c *Cache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.cache[key]; exists {
		delete(c.cache, key)
	}
	return false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for key, value := range c.cache {
		diff := value.time.Second() - time.Now().Second()
		if diff > int(interval.Seconds()) {
			c.Remove(key)
		}
	}

}
