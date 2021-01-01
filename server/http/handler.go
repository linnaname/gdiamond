package main

import (
	"gdiamond/server/service"
	"net/http"
	"net/url"
	"strings"
)

type diamondHandler struct{}

const (
	WORD_SEPARATOR = " "
	LINE_SEPARATOR = "|"
	CONTENT_MD5    = "Content-MD5"
)

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
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	} else {
		w.Write([]byte("Illegal argument,need dataId and group"))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func config(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		dataId := strings.TrimSpace(r.Form.Get("dataId"))
		group := strings.TrimSpace(r.Form.Get("group"))
		if dataId == "" || group == "" {
			goto ARG_ILLEGAL
		} else {
			md5 := service.GetContentMD5(dataId, group)
			if md5 == "" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set(CONTENT_MD5, md5)
			if service.IsModified(dataId, group) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			path := service.GetConfigInfoPath(dataId, group)
			if service.IsModified(dataId, group) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(path))
		}

	} else {
		goto ARG_ILLEGAL
	}

ARG_ILLEGAL:
	w.Write([]byte("Illegal argument,need dataId and group"))
	w.WriteHeader(http.StatusBadRequest)
	return
}

//获取已变更的配置
func getProbeModifyResult(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		probeModify := strings.TrimSpace(r.Form.Get("probeModify"))
		if probeModify == "" {
			goto ARG_ILLEGAL
		} else {
			configKeyList := getConfigKeyList(probeModify)
			resultBuilder := strings.Builder{}
			for i := 0; i < len(configKeyList); i++ {
				dataId := configKeyList[i].DataId
				group := configKeyList[i].Group
				md5 := service.GetContentMD5(dataId, group)
				if md5 != configKeyList[i].MD5 {
					resultBuilder.WriteString(dataId)
					resultBuilder.WriteString(WORD_SEPARATOR)
					resultBuilder.WriteString(group)
					resultBuilder.WriteString(LINE_SEPARATOR)
				}
			}
			result := resultBuilder.String()
			escapeUrl := url.QueryEscape(result)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(escapeUrl))
		}
	} else {
		goto ARG_ILLEGAL
	}

ARG_ILLEGAL:
	w.Write([]byte("Illegal argument,need probeModify"))
	w.WriteHeader(http.StatusBadRequest)
	return
}

type ConfigKey struct {
	DataId string
	Group  string
	MD5    string
}

func getConfigKeyList(probeModify string) []ConfigKey {
	if probeModify == "" {
		return nil
	}
	configKeyStrings := strings.Split(probeModify, LINE_SEPARATOR)
	configKeyList := make([]ConfigKey, len(configKeyStrings))
	for _, configKeyString := range configKeyStrings {
		configKey := strings.Split(configKeyString, WORD_SEPARATOR)
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
