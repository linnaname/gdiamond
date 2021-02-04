package server

import (
	"gdiamond/server/internal/common"
	"gdiamond/server/internal/service"
	"github.com/sirupsen/logrus"
	"log"
)

//Start setup server
func Start() {
	service.SetupLogger()
	err := common.ParseCmdAndInitConfig()
	if err != nil {
		log.Println("ParseCmdAndInitConfig failed", err)
		service.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("ParseCmdAndInitConfig failed")
		return
	}

	err = common.InitDBConn()
	if err != nil {
		log.Println("InitDBConn failed", err)
		service.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("InitDBConn failed")
		return
	}

	err = service.SetupDumpTask()
	if err != nil {
		log.Println("SetupDumpTask", err)
		service.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("SetupDumpTask failed")
		return
	}

	err = service.SetupRegisterTask()
	if err != nil {
		log.Println("SetupRegisterTask", err)
		service.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("SetupRegisterTask failed")
		return
	}
	SetupHttpServer()
}
