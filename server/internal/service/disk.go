package service

import (
	"fmt"
	"gdiamond/server/internal/model"
	"gdiamond/util/fileutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

const configDataDir = "config-data"

var modifyMarkCache sync.Map

//SaveToDisk write config info data to local file
//create dir or file if necessary
func SaveToDisk(info *model.ConfigInfo) error {
	group := info.Group
	dataID := info.DataID
	cacheKey := generateCacheKey(group, dataID)
	_, loaded := modifyMarkCache.LoadOrStore(cacheKey, true)
	if !loaded {
		groupPath := GetFilePath(configDataDir + "/" + group)
		err := fileutil.CreateDirIfNecessary(groupPath)

		if err != nil {
			modifyMarkCache.Delete(cacheKey)
			return err
		}
		targetFile, err := createFileIfNessary(groupPath, dataID)

		if err != nil {
			modifyMarkCache.Delete(cacheKey)
			return err
		}

		tempFile, err := createTempFile(dataID, group)
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
	}
	modifyMarkCache.Delete(cacheKey)
	return fmt.Errorf("config info is being motified, dataID=%s,group=%s", dataID, group)
}

//RemoveConfigInfoFromDisk  remove config info file from disk, it mean delete file
func RemoveConfigInfoFromDisk(dataID, group string) error {
	cacheKey := generateCacheKey(group, dataID)
	_, loaded := modifyMarkCache.LoadOrStore(cacheKey, true)
	if !loaded {
		groupPath := GetFilePath(configDataDir + "/" + group)
		if _, err := os.Stat(groupPath); !os.IsNotExist(err) {
			return nil
		}
		dataPath := GetFilePath(configDataDir + "/" + group + "/" + dataID)
		if _, err := os.Stat(dataPath); !os.IsNotExist(err) {
			return nil
		}
		err := os.Remove(dataPath)
		return err
	}
	modifyMarkCache.Delete(cacheKey)
	return nil
}

//IsModified  whether modified config info by memory cache
func IsModified(dataID, group string) bool {
	v, ok := modifyMarkCache.Load(generateCacheKey(group, dataID))
	if !ok {
		return false
	}
	return v == nil
}

//GetFilePath  filepath where diamond-server store config data
func GetFilePath(dir string) string {
	baseDir := fileutil.GetCurrentDirectory() + "/gdiamond-server"
	return filepath.Join(baseDir, dir)
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

func createTempFile(dataID, group string) (*os.File, error) {
	return ioutil.TempFile("", group+"-"+dataID+".tmp")
}

func generateCacheKey(group, dataID string) string {
	return group + "/" + dataID
}
