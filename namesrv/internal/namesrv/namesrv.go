package namesrv

import (
	"fmt"
	logger "gdiamond/namesrv/internal/log"
	"gdiamond/namesrv/internal/network"
	"gdiamond/namesrv/internal/routeinfo"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type nameServer struct {
	routeInfo *routeinfo.RouteInfo
}

var s *nameServer

//Start  name server
func Start() {
	s = &nameServer{}
	logger.SetupLogger()
	s.setupScanServerTask()
	go setupHttpServer()
	s.startGnetServer()
}

func (s *nameServer) setupScanServerTask() {
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

func (s *nameServer) startGnetServer() {
	logger.Logger.WithFields(logrus.Fields{}).Info("Starting GnetServer")
	//配置化或者使用命令参数
	addr := fmt.Sprintf("tcp://:%d", 9000)
	network.New(addr, s.routeInfo)
}

func setupHttpServer() {
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
			logger.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Close server")
		}
	}()

	log.Println("Starting namesrv  httpserver")
	logger.Logger.WithFields(logrus.Fields{}).Info("Starting namesrv  http server")
	err := server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			logger.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Server closed under request")
		} else {
			logger.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Server closed unexpected")
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
