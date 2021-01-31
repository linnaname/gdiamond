package service

import (
	"errors"
	"gdiamond/server/common"
	"gdiamond/server/model"
	"strings"
	"sync"
	"time"
	"unicode"
)

var cache sync.Map
var locker sync.Mutex

//AddConfigInfo add config to database/save it to disk and notify other gdiamond server nodes
func AddConfigInfo(dataId, group, content string) error {
	err := checkParameter(dataId, group, content)
	if err != nil {
		return err
	}

	configInfo := model.NewConfigInfo(dataId, group, content, time.Now())
	err = addConfingInfo(configInfo)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataId)
	cache.Store(key, configInfo)
	err = SaveToDisk(configInfo)
	if err != nil {
		return err
	}
	NotifyOtherNodes(dataId, group)
	return nil
}

//UpdateConfigInfo update config inof
func UpdateConfigInfo(dataId, group, content string) error {
	err := checkParameter(dataId, group, content)
	if err != nil {
		return err
	}

	configInfo := model.NewConfigInfo(dataId, group, content, time.Now())
	err = updateConfigInfo(configInfo)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataId)
	cache.Store(key, configInfo)
	err = SaveToDisk(configInfo)
	if err != nil {
		return err
	}
	NotifyOtherNodes(dataId, group)
	return nil
}

//LoadConfigInfoToDisk  when other node call NotifyOtherNodes method gdiamond-server will invoke it method
//to load config info from db to disk
func LoadConfigInfoToDisk(dataId, group string) error {
	configInfo, err := findConfigInfo(dataId, group)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataId)
	if configInfo != nil {
		cache.Store(key, configInfo)
		err := SaveToDisk(configInfo)
		if err != nil {
			return err
		}
	} else {
		cache.Delete(key)
		err := RemoveConfigInfoFromDisk(dataId, group)
		if err != nil {
			return err
		}
	}
	return nil
}

//FindConfigInfo find config info from db by dataId and group
func FindConfigInfo(dataId, group string) (*model.ConfigInfo, error) {
	return findConfigInfo(dataId, group)
}

//FindConfigInfoPage find config info by page, group and dataId may be empty
func FindConfigInfoPage(pageNo, pageSize int, group, dataId string) (*model.Page, error) {
	if dataId != "" && group != "" {
		configInfo, err := findConfigInfo(dataId, group)
		if err != nil {
			return nil, err
		}
		page := &model.Page{}
		if configInfo != nil {
			page.PageItems = append(page.PageItems, configInfo)
			page.PageNO = 1
			page.TotalCount = 1
			page.PageAvailable = 1
		}
		return page, nil
	} else if dataId == "" && group != "" {
		return findConfigInfoByGroup(pageNo, pageSize, group)
	} else if dataId != "" && group == "" {
		return findConfigInfoByDataId(pageNo, pageSize, dataId)
	} else {
		return findAllConfigInfo(pageNo, pageSize)
	}
}

func FindConfigInfoLike(pageNo, pageSize int, dataId, group string) (*model.Page, error) {
	return findAllConfigLike(pageNo, pageSize, dataId, group)
}

//NotifyOtherNodes  notify other gdiamond server node when config info changed
func NotifyOtherNodes(dataId, group string) {
	notifyConfigInfoChange(dataId, group)
}

func GetContentMD5(dataId, group string) string {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(dataId, group)
	configInfo, loaded := cache.Load(key)
	if configInfo == nil || !loaded {
		return ""
	}
	value := i2Str(configInfo.(*model.ConfigInfo).MD5)
	return value
}

func GetCache(dataId, group string) *model.ConfigInfo {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(dataId, group)
	value, loaded := cache.Load(key)
	if value == nil || !loaded {
		return nil
	}
	configInfo, _ := value.(*model.ConfigInfo)
	return configInfo
}

func UpdateMD5Cache(configInfo *model.ConfigInfo) {
	key := generateMD5CacheKey(configInfo.DataId, configInfo.Group)
	md5 := common.GetMd5(configInfo.Content)
	configInfo.MD5 = md5
	cache.Store(key, configInfo)
}

//GetConfigInfoPath get local file path by dataId and group
func GetConfigInfoPath(dataId, group string) string {
	builder := strings.Builder{}
	builder.WriteString("/")
	builder.WriteString(configDataDir)
	builder.WriteString("/")
	builder.WriteString(group)
	builder.WriteString("/")
	builder.WriteString(dataId)
	return builder.String()
}

func i2Str(value interface{}) string {
	str, ok := value.(string)
	if ok {
		return str
	}
	return ""
}

func checkParameter(dataId, group, content string) error {
	if dataId == "" || containsWhitespace(dataId) {
		return errors.New("invalid dataId")
	}
	if group == "" || containsWhitespace(group) {
		return errors.New("invalid group")
	}
	if content == "" {
		return errors.New("invalid content")
	}
	return nil
}

func generateMD5CacheKey(dataId, group string) string {
	key := group + "/" + dataId
	return key
}

/**
checks string whether  contains whitespace
*/
func containsWhitespace(token string) bool {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, token) != token
}
