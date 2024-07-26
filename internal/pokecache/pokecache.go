package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	var m sync.Mutex
	cache := Cache{
		cache: make(map[string]cacheEntry),
		mu:    m,
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cache.reapLoop(interval)
			}
			return
		}

	}()

	return cache
}

func (cache *Cache) reapLoop(interval time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	currentTime := time.Now()
	for cachekey, cacheVal := range cache.cache {
		cacheExpiration := cacheVal.createdAt
		cacheExpiration = cacheExpiration.Add(interval)

		// If the current time is after the expiration time, delete the entry
		if currentTime.After(cacheExpiration) {
			delete(cache.cache, cachekey)
		}
	}
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	data := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.cache[key] = data
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if value, exists := cache.cache[key]; exists {
		return value.val, true
	} else {
		return nil, false
	}
}
