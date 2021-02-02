package main

import (
	"fmt"
	"gdiamond/namesrv/network"
	"gdiamond/namesrv/routeinfo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type startUp struct {
	routeInfo *routeinfo.RouteInfo
}

var s *startUp

func main() {
	s = &startUp{}
	s.initialize()
	go httpServer()
	s.start()
}

func (s *startUp) initialize() {
	s.routeInfo = routeinfo.New()

	scanTicker := time.NewTicker(time.Second * 60)
	go func() {
		defer scanTicker.Stop()
		for {
			<-scanTicker.C
			s.routeInfo.ScanNotActiveServer()
		}
	}()
}

func (s *startUp) start() {
	//配置化或者使用命令参数
	addr := fmt.Sprintf("tcp://:%d", 9000)
	network.New(addr, s.routeInfo)
}

func httpServer() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	mux := http.NewServeMux()
	mux.HandleFunc("/namesrv/addrs", cluster)

	server := &http.Server{
		Addr:         ":8080",
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

func cluster(w http.ResponseWriter, r *http.Request) {
	cInfo := s.routeInfo.GetAllClusterInfo()
	builder := strings.Builder{}
	for _, v := range cInfo.ServerAddrTable {
		serverAddres := v.GetServerAddrs()
		if len(serverAddres) > 0 {
			for _, serverAddr := range serverAddres {
				builder.WriteString(serverAddr)
				builder.WriteString("\n")
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(builder.String()))
}
