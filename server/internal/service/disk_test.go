package service

import (
	"fmt"
	"gdiamond/server/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveToDisk(t *testing.T) {
	configInfo := &model.ConfigInfo{Group: "DEFAULT_GROUP", DataID: "linna.com", Content: "song for linana", MD5: "6f6e0326c63ed62c709d874f7093f6e1"}
	err := SaveToDisk(configInfo)
	assert.NoError(t, err)
}

func TestIsModified(t *testing.T) {
	v := IsModified("DEFAULT_GROUP", "linna.com")
	assert.False(t, v)
}

func TestRemoveConfigInfoFromDisk(t *testing.T) {
	err := RemoveConfigInfoFromDisk("linna.com", "DEFAULT_GROUP")
	assert.NoError(t, err)
}

func TestGetFilePath(t *testing.T) {
	assert.NotEmpty(t, GetFilePath("group/dataID"))
	fmt.Println(GetFilePath(""))
}
