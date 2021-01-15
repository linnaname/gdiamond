package processor

import (
	"errors"
	"gdiamond/util/fileutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SnapshotConfigInfoProcessor struct {
	dir string
}

func NewSnapshotConfigInfoProcessor(dir string) *SnapshotConfigInfoProcessor {
	processor := &SnapshotConfigInfoProcessor{
		dir: dir,
	}
	fileutil.CreateDirIfNessary(dir)
	return processor
}

/**
从文件中读取snapshot
*/
func (p *SnapshotConfigInfoProcessor) GetConfigInfomation(dataId, group string) (string, error) {
	if dataId == "" || group == "" {
		return "", nil
	}

	path := filepath.Join(p.dir, group)
	if !fileutil.IsExist(path) {
		return "", nil
	}
	filePath := filepath.Join(path, dataId)
	if !fileutil.IsExist(filePath) {
		return "", nil
	}

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

/**
保存snapshot
*/
func (p *SnapshotConfigInfoProcessor) SaveSnaptshot(dataId, group, config string) error {
	if dataId == "" {
		errors.New("blank dataId")
	}
	if group == "" {
		errors.New("blank group")
	}

	if config == "" {
		config = ""
	}

	f, err := p.getTargetFile(dataId, group)
	if err != nil {
		return err
	}
	writeN, err := io.WriteString(f, config)
	if writeN == 0 || err != nil {
		return err
	}
	return nil
}

/**
 * 删除snapshot
 */
func (p *SnapshotConfigInfoProcessor) RemoveSnapshot(dataId, group string) {
	if dataId == "" || group == "" {
		return
	}

	path := filepath.Join(p.dir, group)
	if !fileutil.IsExist(path) {
		return
	}

	filePath := filepath.Join(path, dataId)
	if !fileutil.IsExist(filePath) {
		return
	}
	os.Remove(filePath)

	// 如果目录没有文件了，删除目录
	list, _ := ioutil.ReadDir(path)
	if list == nil || len(list) == 0 {
		os.Remove(path)
	}
}

func (p *SnapshotConfigInfoProcessor) getTargetFile(dataId, group string) (*os.File, error) {
	path := filepath.Join(p.dir, group)
	fileutil.CreateDirIfNessary(path)
	filePath := filepath.Join(path, dataId)
	return fileutil.CreateFileIfNessary(filePath)
}
