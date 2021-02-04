package server

import (
	"gdiamond/server/internal/common"
	"gdiamond/server/internal/service"
	"gdiamond/util/fileutil"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
)

const (
	logDir  = "logs"
	logFile = "gdiamond.log"
)

var Logger = logrus.New()

//Start setup server
func Start() {
	setupLogger()
	err := common.ParseCmdAndInitConfig()
	if err != nil {
		log.Println("ParseCmdAndInitConfig failed", err)
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("ParseCmdAndInitConfig failed")
		return
	}

	err = common.InitDBConn()
	if err != nil {
		log.Println("InitDBConn failed", err)
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("InitDBConn failed")
		return
	}

	err = service.SetupDumpTask()
	if err != nil {
		log.Println("SetupDumpTask", err)
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("SetupDumpTask failed")
		return
	}

	err = service.SetupRegisterTask()
	if err != nil {
		log.Println("SetupRegisterTask", err)
		Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("SetupRegisterTask failed")
		return
	}
	SetupHttpServer()
}

func setupLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	filePath := service.GetFilePath(logDir)
	fileutil.CreateDirIfNecessary(filePath)
	file, err := os.OpenFile(filepath.Join(filePath, logFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = file
	} else {
		Logger.Warn("Failed to log to file, using default stderr")
	}
}
