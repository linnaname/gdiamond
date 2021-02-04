package service

import (
	"gdiamond/util/fileutil"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const (
	logDir  = "logs"
	logFile = "gdiamond.log"
)

var Logger = logrus.New()

//SetupLogger setup logrus path and format
func SetupLogger() {
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	filePath := GetFilePath(logDir)
	fileutil.CreateDirIfNecessary(filePath)
	file, err := os.OpenFile(filepath.Join(filePath, logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = file
	} else {
		Logger.Warn("Failed to log to file, using default stderr")
	}
}
