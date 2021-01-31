package service

import (
	"errors"
	"fmt"
	"gdiamond/server/model"
	"gdiamond/util/fileutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const configDataDir = "config-data"

var modifyMarkCache sync.Map

func SaveToDisk(info *model.ConfigInfo) error {
	group := info.Group
	dataId := info.DataId
	cacheKey := generateCacheKey(group, dataId)
	_, loaded := modifyMarkCache.LoadOrStore(cacheKey, true)
	if !loaded {
		groupPath := GetFilePath(configDataDir + "/" + group)
		err := fileutil.CreateDirIfNessary(groupPath)

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

func RemoveConfigInfoFromDisk(dataId, group string) error {
	cacheKey := generateCacheKey(group, dataId)
	_, loaded := modifyMarkCache.LoadOrStore(cacheKey, true)
	if !loaded {
		groupPath := GetFilePath(configDataDir + "/" + group)
		if _, err := os.Stat(groupPath); !os.IsNotExist(err) {
			return nil
		}
		dataPath := GetFilePath(configDataDir + "/" + group + "/" + dataId)
		if _, err := os.Stat(dataPath); !os.IsNotExist(err) {
			return nil
		}
		err := os.Remove(dataPath)
		return err
	}
	modifyMarkCache.Delete(cacheKey)
	return nil
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
	deleteTempFile(tempFile)
}

func deleteTempFile(tempFile *os.File) {
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		os.Remove(tempFile.Name())
	}
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

func GetFilePath(dir string) string {
	baseDir := fileutil.GetCurrentDirectory() + "/gdiamond-server"
	return filepath.Join(baseDir, dir)
}
