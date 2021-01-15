package configinfo

import (
	"sync"
	"time"
)

type SimpleCache struct {
	cacheTTL int64
	cache    sync.Map
}

type CacheEntry struct {
	timestamp int64
	value     interface{}
}

func NewSimpleCache(cacheTTL int64) *SimpleCache {
	cache := &SimpleCache{cacheTTL: cacheTTL}
	return cache
}

func (s *SimpleCache) put(key string, e interface{}) {
	if key == "" || e == nil {
		return
	}
	entry := CacheEntry{timestamp: time.Now().UnixNano() + s.cacheTTL, value: e}
	s.cache.Store(key, entry)
}

func (s *SimpleCache) get(key string) interface{} {
	var result interface{}
	value, ok := s.cache.Load(key)
	if ok && value != nil {
		entry := value.(CacheEntry)
		if entry.timestamp > time.Now().UnixNano() {
			result = entry.value
		}
	}
	return result
}
