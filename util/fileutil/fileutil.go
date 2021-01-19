package fileutil

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

/**
CreateDirIfNessary create directory if not exist,pay attention to fix permission and
it will create any necessary parents
If path is already a exist, it does nothing and returns nil.
*/
func CreateDirIfNessary(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

//IsExist if file or directory exist,don't care about err
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsExist(err) || err == nil {
		return true
	}
	return false
}

/**
CreateFileIfNessary create file if not exist,pay attention to fix permission
If file is already a exist, it does nothing and returns nil.
*/
func CreateFileIfNessary(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filepath)
		return file, err
	}
	return file, err
}

/**
GetFileContent get all file content once
*/
func GetFileContent(filePath string) (string, error) {
	finfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	if finfo.IsDir() {
		return "", errors.New("Not file")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return "", errors.New("Can't open file")
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.New("Can't read file")
	}
	return string(content), nil
}

/**
GetGrandpaDir get grandpa base dir of given filepath
path must be file path not dir
return  grandpa dir path,if there no grandpa it will return empty string and error
*/
func GetGrandpaDir(path string) (string, error) {
	if IsDir(path) {
		return "", errors.New("not valid file")
	}
	parentPath := filepath.Dir(path)
	if IsDir(parentPath) {
		grandpaPath := filepath.Dir(parentPath)
		if IsDir(grandpaPath) {
			return filepath.Base(grandpaPath), nil
		} else {
			return "", errors.New("grandpa path not dir")
		}
	} else {
		return "", errors.New("parent path not dir")
	}
}

//IsDir if path is dir,don't care about err
func IsDir(path string) bool {
	finfo, _ := os.Stat(path)
	return finfo.IsDir()
}
