package simplecache

import (
	"sync"
	"time"
)

//SimpleCache simple cache with ttl
type SimpleCache struct {
	//ttl for whole cache,not for single key, time unit:second
	ttl   int64
	cache sync.Map
}

//CacheEntry entry
type CacheEntry struct {
	//exipred timestamp
	timestamp int64
	value     interface{}
}

//New new
func New(ttl int64) *SimpleCache {
	cache := &SimpleCache{ttl: ttl}
	return cache
}

//Put put k,v to cache
func (s *SimpleCache) Put(key string, e interface{}) {
	if key == "" || e == nil {
		return
	}
	entry := CacheEntry{timestamp: time.Now().Unix() + s.ttl, value: e}
	s.cache.Store(key, entry)
}

//Get get value of key
//if not exist return nil interface
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
