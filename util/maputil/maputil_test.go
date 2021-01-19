package maputil

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLength(t *testing.T) {
	var m sync.Map
	assert.NotNil(t, m)
	l := LengthOfSyncMap(m)
	assert.Equal(t, int64(0), l)
	m.Store("test", "tv")
	assert.Equal(t, int64(1), LengthOfSyncMap(m))
}

func TestClearSyncMap(t *testing.T) {
	var m sync.Map
	assert.NotNil(t, m)
	l := LengthOfSyncMap(m)
	assert.Equal(t, int64(0), l)
	m.Store("test", "tv")
	assert.Equal(t, int64(1), LengthOfSyncMap(m))
	ClearSyncMap(m)
	assert.NotNil(t, m)
	assert.Equal(t, int64(0), l)
}
