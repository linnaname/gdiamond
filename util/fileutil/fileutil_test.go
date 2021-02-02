package fileutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestDir  = "test"
	TestFile = "test.file"
)

func TestCreateDirIfNecessary(t *testing.T) {
	err := CreateDirIfNecessary(TestDir)
	assert.NoError(t, err)
}

func TestIsDir(t *testing.T) {
	assert.True(t, IsDir(TestDir))
}

func TestIsExist(t *testing.T) {
	b := IsExist(TestDir)
	assert.True(t, b)
}

func TestCreateFileIfNessary(t *testing.T) {
	f, err := CreateFileIfNessary(TestFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	assert.True(t, IsExist(f.Name()))
	assert.False(t, IsDir(f.Name()))
}

func TestGetFileContent(t *testing.T) {
	content, err := GetFileContent(TestFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestGetGrandpaDir(t *testing.T) {
	g1, err := GetGrandpaDir(TestDir)
	assert.Error(t, err)
	assert.Empty(t, g1)

	g2, err := GetGrandpaDir(TestFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, g2)
}

func TestString2File(t *testing.T) {
	err := String2File("11", "./test/kv.json")
	assert.NoError(t, err)
}

func TestMMapRead(t *testing.T) {
	buf, err := MMapRead(TestFile)
	assert.NoError(t, err)
	assert.NotNil(t, buf)
}

func TestGetCurrentDirectory(t *testing.T) {
	assert.NotEmpty(t, GetCurrentDirectory())
	fmt.Println(GetCurrentDirectory())
}
