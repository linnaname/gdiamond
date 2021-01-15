package simplecache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	cache := New(10)
	value := cache.get("test")
	assert.Nil(t, value)
	cache.put("test", "val")
	assert.NotNil(t, cache.get("test"))
}

func TestSimpleCacheTime(t *testing.T) {
	cache := New(10)
	value := cache.get("test")
	assert.Nil(t, value)
	cache.put("test", "val")
	assert.NotNil(t, cache.get("test"))
	time.Sleep(time.Second * 12)
	assert.Nil(t, cache.get("test"))
}
