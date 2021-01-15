package fileutil

import "os"

/**
目录不存在则创建
*/
func CreateDirIfNessary(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsExist(err) || err == nil {
		return true
	}
	return false
}

func CreateFileIfNessary(filepath string) (*os.File, error) {
	file, err := os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filepath)
		return file, err
	}
	return file, err
}
