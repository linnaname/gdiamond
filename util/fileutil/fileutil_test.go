package fileutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsExist(t *testing.T) {
	b := IsExist("/Users/goranka/tmp/test/test.group")
	assert.True(t, b)
}

func TestGetFileContent(t *testing.T) {
	content, err := GetFileContent("/Users/goranka/tmp/test/e.dir/e.data")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestGetGrandpaDir(t *testing.T) {
	g1, err := GetGrandpaDir("/Users/goranka/tmp/test/e.dir")
	assert.Error(t, err)
	assert.Empty(t, g1)

	g2, err := GetGrandpaDir("/Users/goranka/tmp/test/e.dir/e.data")
	assert.NoError(t, err)
	assert.Equal(t, "/Users/goranka/tmp/test", g2)
}
