package server

import (
	"gdiamond/server/internal/service"
	"gdiamond/util/fileutil"
	"gdiamond/util/netutil"
	"gdiamond/util/stringutil"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type diamondHandler struct{}

const (
	wordSeparator            = ","
	lineSeparator            = ";"
	contentMd5               = "Content-MD5"
	lastModified             = "Last-Modified"
	probeModifyRequest       = "Probe-Modify-Request"
	longPollingTimeOutHeader = "Long-Polling-TimeOut"
)

var notifier chan *Event

type Event struct {
	dataId    string
	group     string
	timestamp int64
}

func SetupHttpServer() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	notifier = make(chan *Event)

	mux := http.NewServeMux()
	mux.Handle("/", &diamondHandler{})
	mux.HandleFunc("/diamond-server/notify", notifyConfigInfo)
	mux.HandleFunc("/diamond-server/config", config)
	mux.HandleFunc("/diamond-server/getProbeModify", getProbeModify)
	mux.HandleFunc("/diamond-server/publishConfig", publishConfig)

	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 90,
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
	service.Logger.WithFields(logrus.Fields{}).Info("Starting  httpserver")

	err := server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			service.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Server closed under request")
		} else {
			service.Logger.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Fatal("Server closed unexpected")
		}
	}
}

func (*diamondHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This diamond http server"))
}

//notifyConfigInfo http method
func notifyConfigInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		dataId := strings.TrimSpace(r.Form.Get("dataId"))
		group := strings.TrimSpace(r.Form.Get("group"))
		err := service.LoadConfigInfoToDisk(dataId, group)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Load config to disk successed"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Illegal argument,need dataId and group"))
}

func config(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		dataId := strings.TrimSpace(r.Form.Get("dataId"))
		group := strings.TrimSpace(r.Form.Get("group"))
		if dataId == "" || group == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Illegal argument,need dataId and group"))
			return
		} else {
			cacheInfo := service.GetCache(dataId, group)
			if cacheInfo == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set(contentMd5, cacheInfo.MD5)
			w.Header().Set(lastModified, cacheInfo.LastModified.String())
			if service.IsModified(dataId, group) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			if service.IsModified(dataId, group) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			path := service.GetConfigInfoPath(dataId, group)
			filePath := service.GetFilePath(path)
			buf, err := fileutil.MMapRead(filePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(buf)
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal argument,need dataId and group"))
		return
	}
}

func publishConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		dataId := strings.TrimSpace(r.Form.Get("dataId"))
		group := strings.TrimSpace(r.Form.Get("group"))
		content := strings.TrimSpace(r.Form.Get("content"))
		errorMessage := "illegal argument"
		checkSuccess := true
		if stringutil.HasInvalidChar(dataId) {
			checkSuccess = false
			errorMessage = "invalid DataID"
		}
		if stringutil.HasInvalidChar(group) {
			checkSuccess = false
			errorMessage = "invalid group"
		}

		if stringutil.HasInvalidChar(content) {
			checkSuccess = false
			errorMessage = "invalid content"
		}
		if !checkSuccess {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errorMessage))
			return
		}

		clientIp := netutil.GetRemoteClientIP(r)
		if clientIp == "" {
			service.Logger.WithFields(logrus.Fields{
				"request": r,
				"dataId":  dataId,
				"group":   group,
			}).Error("clientIp is empty")
		}

		err := service.AddOrUpdate(dataId, group, content)
		if err != nil {
			service.Logger.WithFields(logrus.Fields{
				"content": content,
				"dataId":  dataId,
				"group":   group,
				"err":     err.Error(),
			}).Error("service.AddOrUpdate failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//暂时先强入侵业务逻辑，后面看有没有更好的办法
		go notifyListener(dataId, group)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal argument"))
		return
	}
}

func notifyListener(dataId, group string) {
	event := &Event{dataId: dataId, group: group, timestamp: time.Now().Unix()}
	notifier <- event
}

//获取已变更的配置
func getProbeModify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		probeModify := strings.TrimSpace(r.Form.Get(probeModifyRequest))
		if probeModify == "" {
			goto ARG_ILLEGAL
		} else {
			//modifyResult := getModify(probeModify)
			//if modifyResult != "" {
			//	escapeURL := url.QueryEscape(modifyResult)
			//	service.Logger.WithFields(logrus.Fields{}).Debug("getModify in")
			//	w.WriteHeader(http.StatusOK)
			//	w.Write([]byte(escapeURL))
			//	return
			//}

			lptHeader := r.Header.Get(longPollingTimeOutHeader)
			longPollingTimeout, _ := strconv.Atoi(lptHeader)
			ctx := r.Context()

			select {
			case <-notifier:
				modifyResult := getModify(probeModify)
				if modifyResult != "" {
					escapeURL := url.QueryEscape(modifyResult)
					service.Logger.WithFields(logrus.Fields{}).Debug("event notifier")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(escapeURL))
					return
				}
			case <-time.After(time.Millisecond * time.Duration(longPollingTimeout)):
				service.Logger.WithFields(logrus.Fields{}).Debug("hangup time out")
				w.WriteHeader(http.StatusNotModified)
				return
			case <-ctx.Done():
				service.Logger.WithFields(logrus.Fields{}).Debug("Client has disconnected")
				w.WriteHeader(http.StatusOK)
				return
			}

		}
	} else {
		goto ARG_ILLEGAL
	}

ARG_ILLEGAL:
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Illegal argument,need probeModify"))
	return
}

func getModify(probeModify string) string {
	configKeyList := getConfigKeyList(probeModify)
	resultBuilder := strings.Builder{}
	for i := 0; i < len(configKeyList); i++ {
		dataId := configKeyList[i].DataId
		group := configKeyList[i].Group
		md5 := service.GetContentMD5(dataId, group)
		if md5 != configKeyList[i].MD5 {
			resultBuilder.WriteString(dataId)
			resultBuilder.WriteString(wordSeparator)
			resultBuilder.WriteString(group)
			resultBuilder.WriteString(lineSeparator)
		}
	}
	return resultBuilder.String()
}

//ConfigKey to simplify config model op
type ConfigKey struct {
	DataId string
	Group  string
	MD5    string
}

func getConfigKeyList(probeModify string) []ConfigKey {
	if probeModify == "" {
		return nil
	}
	configKeyStrings := strings.Split(probeModify, lineSeparator)
	configKeyList := make([]ConfigKey, 0)
	for _, configKeyString := range configKeyStrings {
		configKey := strings.Split(configKeyString, wordSeparator)
		if len(configKey) > 3 {
			continue
		}
		if configKey[0] == "" {
			continue
		}
		key := ConfigKey{}
		key.DataId = configKey[0]
		if len(configKey) >= 2 && configKey[1] != "" {
			key.Group = configKey[1]
		}
		if len(configKey) == 3 && configKey[2] != "" {
			key.MD5 = configKey[2]
		}
		configKeyList = append(configKeyList, key)
	}
	return configKeyList
}
