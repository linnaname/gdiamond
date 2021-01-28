package fileutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TEST_DIR  = "test"
	TEST_FILE = "test.file"
)

func TestCreateDirIfNessary(t *testing.T) {
	err := CreateDirIfNessary(TEST_DIR)
	assert.NoError(t, err)
}

func TestIsDir(t *testing.T) {
	assert.True(t, IsDir(TEST_DIR))
}

func TestIsExist(t *testing.T) {
	b := IsExist(TEST_DIR)
	assert.True(t, b)
}

func TestCreateFileIfNessary(t *testing.T) {
	f, err := CreateFileIfNessary(TEST_FILE)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.True(t, IsExist(f.Name()))
	assert.False(t, IsDir(f.Name()))
}

func TestGetFileContent(t *testing.T) {
	content, err := GetFileContent(TEST_FILE)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestGetGrandpaDir(t *testing.T) {
	g1, err := GetGrandpaDir(TEST_DIR)
	assert.Error(t, err)
	assert.Empty(t, g1)

	g2, err := GetGrandpaDir(TEST_FILE)
	assert.NoError(t, err)
	assert.NotEmpty(t, g2)
}

func TestString2File(t *testing.T) {
	err := String2File("11", "./test/kv.json")
	assert.NoError(t, err)
}

func TestMMapRead(t *testing.T) {
	buf, err := MMapRead(TEST_FILE)
	assert.NoError(t, err)
	assert.NotNil(t, buf)
}
