package server

import (
	"fmt"
	"gdiamond/server/internal/common"
	"gdiamond/server/internal/service"
)

//Start setup server
func Start() {
	//TODO err handle
	err := common.ParseCmdAndInitConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = common.InitDBConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = service.SetupDumpTask()
	if err != nil {
		fmt.Println(err)
		return
	}
	service.SetupRegisterTask()
	SetupHttpServer()
}
