package service

import (
	"errors"
	"fmt"
	"gdiamond/server/model"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const BASE_DIR = "config-data"

var modifyMarkCache sync.Map

func SaveToDisk(info *model.ConfigInfo) error {
	group := info.Group
	dataId := info.DataId
	cacheKey := generateCacheKey(group, dataId)
	_, loaded := modifyMarkCache.LoadOrStore(cacheKey, true)
	if !loaded {
		groupPath := getFilePath(BASE_DIR + "/" + group)
		err := createDirIfNessary(groupPath)

		if err != nil {
			modifyMarkCache.Delete(cacheKey)
			return err
		}
		targetFile, err := createFileIfNessary(groupPath, dataId)

		if err != nil {
			modifyMarkCache.Delete(cacheKey)
			return err
		}

		tempFile, err := createTempFile(dataId, group)
		if err != nil {
			modifyMarkCache.Delete(cacheKey)
			return err
		}

		writeN, err := io.WriteString(tempFile, info.Content)
		if writeN == 0 || err != nil {
			clearCacheAndFile(tempFile, cacheKey)
			return err
		}

		//this is important,switch write to read
		tempFile.Seek(0, 0)
		_, cerr := io.Copy(targetFile, tempFile)
		if cerr != nil {
			clearCacheAndFile(tempFile, cacheKey)
			return cerr
		}

		clearCacheAndFile(tempFile, cacheKey)
		return nil
	} else {
		modifyMarkCache.Delete(cacheKey)
		return errors.New(fmt.Sprintf("config info is being motified, dataId=%s,group=%s", dataId, group))
	}
}

func IsModified(dataId, group string) bool {
	v, ok := modifyMarkCache.Load(generateCacheKey(group, dataId))
	if !ok {
		return false
	}
	return v == nil
}

func clearCacheAndFile(tempFile *os.File, cacheKey string) {
	modifyMarkCache.Delete(cacheKey)
	deteTempFile(tempFile)
}

func deteTempFile(tempFile *os.File) {
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		os.Remove(tempFile.Name())
	}
}

func createDirIfNessary(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func createFileIfNessary(parent, child string) (*os.File, error) {
	name := filepath.Join(parent, child)
	file, err := os.OpenFile(name, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(name)
		return file, err
	}
	return file, err
}

func createTempFile(dataId, group string) (*os.File, error) {
	return ioutil.TempFile("", group+"-"+dataId+".tmp")
}

func generateCacheKey(group, dataId string) string {
	return group + "/" + dataId
}

func getFilePath(dir string) string {
	return filepath.Join(getCurrentDirectory(), dir)
}

func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
