package simplecache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimpleCache(t *testing.T) {
	cache := New(10)
	value := cache.Get("test")
	assert.Nil(t, value)
	cache.Put("test", "val")
	assert.NotNil(t, cache.Get("test"))
}

func TestSimpleCacheTime(t *testing.T) {
	cache := New(10)
	value := cache.Get("test")
	assert.Nil(t, value)
	cache.Put("test", "val")
	assert.NotNil(t, cache.Get("test"))
	time.Sleep(time.Second * 12)
	assert.Nil(t, cache.Get("test"))
}
