package simplecache

import (
	"sync"
	"time"
)

type SimpleCache struct {
	cacheTTL int64 //秒
	cache    sync.Map
}

type CacheEntry struct {
	timestamp int64
	value     interface{}
}

func New(cacheTTL int64) *SimpleCache {
	cache := &SimpleCache{cacheTTL: cacheTTL}
	return cache
}

func (s *SimpleCache) Put(key string, e interface{}) {
	if key == "" || e == nil {
		return
	}
	entry := CacheEntry{timestamp: time.Now().Unix() + s.cacheTTL, value: e}
	s.cache.Store(key, entry)
}

func (s *SimpleCache) Get(key string) interface{} {
	var result interface{}
	value, ok := s.cache.Load(key)
	if ok && value != nil {
		entry := value.(CacheEntry)
		if entry.timestamp > time.Now().Unix() {
			result = entry.value
		}
	}
	return result
}
