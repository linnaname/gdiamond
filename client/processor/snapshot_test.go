package processor

import (
	"gdiamond/util/fileutil"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

const SNAPSHOT_PATH = "/Users/goranka/tmp/test"

func TestSnapshotConfigInfoProcessor_SaveSnaptshot(t *testing.T) {
	p := NewSnapshotConfigInfoProcessor(SNAPSHOT_PATH)
	err := p.SaveSnaptshot("test.dataId", "test.group", "this is content")
	assert.NoError(t, err)
}

func TestSnapshotConfigInfoProcessor_GetConfigInfomation(t *testing.T) {
	p := NewSnapshotConfigInfoProcessor(SNAPSHOT_PATH)
	content, err := p.GetConfigInfomation("test.dataId", "test.group")
	assert.NoError(t, err)
	assert.Equal(t, content, "this is content")
}

func TestSnapshotConfigInfoProcessor_RemoveSnapshot(t *testing.T) {
	p := NewSnapshotConfigInfoProcessor(SNAPSHOT_PATH)
	p.RemoveSnapshot("test.dataId", "test.group")
	assert.False(t, fileutil.IsExist(filepath.Join(SNAPSHOT_PATH, "test.dataId", "test.group")))
}
