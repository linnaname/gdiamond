package fileutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsExist(t *testing.T) {
	b := IsExist("/Users/goranka/tmp/test/test.group")
	assert.True(t, b)
}
