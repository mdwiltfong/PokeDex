package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	time time.Time
	val  []byte
}

type Cache struct {
	cache map[string]CacheEntry
	mu    *sync.Mutex
}

func (cache *Cache) Length() int {
	return len(cache.cache)
}
func NewCache(interval time.Duration) *Cache {
	c := Cache{
		cache: make(map[string]CacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CacheEntry{
		time: time.Now(),
		val:  val,
	}
	fmt.Println("Cache Add: ", c.cache, key, val)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.cache[key]
	fmt.Println("Cache Get: ", c.cache, key, value, exists)
	if exists {
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
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cache {
		if v.time.Before(now.Add(-last)) {
			delete(c.cache, k)
		}
	}
}
