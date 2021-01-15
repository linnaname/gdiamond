package configinfo

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMapNil(t *testing.T) {
	var c sync.Map
	assert.Nil(t, c)
}
