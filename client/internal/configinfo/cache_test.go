package configinfo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCacheData(t *testing.T) {
	c := NewCacheData("linname", " DEFAULT_GROUP")
	assert.NotNil(t, c)
}

func TestCacheData_IncrementFetchCountAndGet(t *testing.T) {
	c := NewCacheData("linname", " DEFAULT_GROUP")
	assert.NotNil(t, c)
	assert.Equal(t, c.GetFetchCount(), int64(0))
	assert.Equal(t, c.IncrementFetchCountAndGet(), int64(1))
	assert.Equal(t, c.GetFetchCount(), int64(1))
}
