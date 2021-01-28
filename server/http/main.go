package main

import (
	"gdiamond/server/common"
	"gdiamond/server/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var server *http.Server

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	common.InitConfig()
	common.InitDBConn()
	service.DumpAll2Disk()
	service.Init()

	register := &service.Register{}
	register.RegisterServerAll()
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			register.RegisterServerAll()
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", &diamondHandler{})
	mux.HandleFunc("/diamond-server/notify", notifyConfigInfo)
	mux.HandleFunc("/diamond-server/config", config)
	mux.HandleFunc("/diamond-server/getProbeModify", getProbeModifyResult)

	server = &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}

	go func() {
		// 接收退出信号
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal("Close server:", err)
		}
	}()

	log.Println("Starting  httpserver")
	err := server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			log.Fatal("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected", err)
		}
	}
}
