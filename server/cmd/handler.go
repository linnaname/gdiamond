package main

import (
	"gdiamond/server/service"
	"gdiamond/util/fileutil"
	"gdiamond/util/stringutil"
	"net/http"
	"net/url"
	"strings"
)

type diamondHandler struct{}

const (
	wordSeparator      = ","
	lineSeparator      = ";"
	contentMd5         = "Content-MD5"
	lastModified       = "Last-Modified"
	probeModifyRequest = "Probe-Modify-Request"
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
		err := service.AddOrUpdate(dataId, group, content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal argument"))
		return
	}
}

//获取已变更的配置
func getProbeModifyResult(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) > 0 {
		probeModify := strings.TrimSpace(r.Form.Get(probeModifyRequest))
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
					resultBuilder.WriteString(wordSeparator)
					resultBuilder.WriteString(group)
					resultBuilder.WriteString(lineSeparator)
				}
			}
			result := resultBuilder.String()
			escapeURL := url.QueryEscape(result)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(escapeURL))
			return
		}
	} else {
		goto ARG_ILLEGAL
	}

ARG_ILLEGAL:
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Illegal argument,need probeModify"))
	return
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
