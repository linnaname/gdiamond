package service

import (
	"errors"
	"gdiamond/server/internal/model"
	"gdiamond/util/stringutil"
	"strings"
	"sync"
	"time"
	"unicode"
)

var cache sync.Map
var locker sync.Mutex

func AddOrUpdate(dataID, group, content string) error {
	cInfo, err := findConfigInfo(dataID, group)
	if err != nil {
		return err
	}
	if cInfo != nil {
		return UpdateConfigInfo(dataID, group, content)
	}
	return AddConfigInfo(dataID, group, content)
}

//AddConfigInfo add config to database/save it to disk and notify other gdiamond server nodes
func AddConfigInfo(dataID, group, content string) error {
	err := checkParameter(dataID, group, content)
	if err != nil {
		return err
	}

	configInfo := model.NewConfigInfo(dataID, group, content, time.Now())
	err = addConfigInfo(configInfo)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataID)
	cache.Store(key, configInfo)
	err = SaveToDisk(configInfo)
	if err != nil {
		return err
	}
	NotifyOtherNodes(dataID, group)
	return nil
}

//UpdateConfigInfo update config info
func UpdateConfigInfo(dataID, group, content string) error {
	err := checkParameter(dataID, group, content)
	if err != nil {
		return err
	}

	configInfo := model.NewConfigInfo(dataID, group, content, time.Now())
	err = updateConfigInfo(configInfo)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataID)
	cache.Store(key, configInfo)
	err = SaveToDisk(configInfo)
	if err != nil {
		return err
	}
	NotifyOtherNodes(dataID, group)
	return nil
}

//LoadConfigInfoToDisk  when other node call NotifyOtherNodes method gdiamond-server will invoke it method
//to load config info from db to disk
func LoadConfigInfoToDisk(dataID, group string) error {
	configInfo, err := findConfigInfo(dataID, group)
	if err != nil {
		return err
	}
	key := generateCacheKey(group, dataID)
	if configInfo != nil {
		cache.Store(key, configInfo)
		err := SaveToDisk(configInfo)
		if err != nil {
			return err
		}
	} else {
		cache.Delete(key)
		err := RemoveConfigInfoFromDisk(dataID, group)
		if err != nil {
			return err
		}
	}
	return nil
}

//FindConfigInfo find config info from db by dataID and group
func FindConfigInfo(dataID, group string) (*model.ConfigInfo, error) {
	return findConfigInfo(dataID, group)
}

//FindConfigInfoPage find config info by page, group and dataID may be empty
func FindConfigInfoPage(pageNo, pageSize int, group, dataID string) (*model.Page, error) {
	if dataID != "" && group != "" {
		configInfo, err := findConfigInfo(dataID, group)
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
	} else if dataID == "" && group != "" {
		return findConfigInfoByGroup(pageNo, pageSize, group)
	} else if dataID != "" && group == "" {
		return findConfigInfoByDataID(pageNo, pageSize, dataID)
	} else {
		return findAllConfigInfo(pageNo, pageSize)
	}
}

//FindConfigInfoLike find config info by like dataID and group
func FindConfigInfoLike(pageNo, pageSize int, dataID, group string) (*model.Page, error) {
	return findAllConfigLike(pageNo, pageSize, dataID, group)
}

//NotifyOtherNodes  notify other gdiamond server node when config info changed
func NotifyOtherNodes(dataID, group string) {
	notifyConfigInfoChange(dataID, group)
}

//GetContentMD5 get memory cache md5 from dataID and group
func GetContentMD5(dataID, group string) string {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(dataID, group)
	configInfo, loaded := cache.Load(key)
	if configInfo == nil || !loaded {
		return ""
	}
	value := i2Str(configInfo.(*model.ConfigInfo).MD5)
	return value
}

//GetCache get memory cache by dataID and group
func GetCache(dataID, group string) *model.ConfigInfo {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(dataID, group)
	value, loaded := cache.Load(key)
	if value == nil || !loaded {
		return nil
	}
	configInfo, _ := value.(*model.ConfigInfo)
	return configInfo
}

//UpdateMD5Cache update memory cache, the most important is md5
func UpdateMD5Cache(configInfo *model.ConfigInfo) {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(configInfo.DataID, configInfo.Group)
	md5 := stringutil.GetMd5(configInfo.Content)
	configInfo.MD5 = md5
	cache.Store(key, configInfo)
}

//GetConfigInfoPath get local file path by dataID and group
func GetConfigInfoPath(dataID, group string) string {
	builder := strings.Builder{}
	builder.WriteString("/")
	builder.WriteString(configDataDir)
	builder.WriteString("/")
	builder.WriteString(group)
	builder.WriteString("/")
	builder.WriteString(dataID)
	return builder.String()
}

func i2Str(value interface{}) string {
	str, ok := value.(string)
	if ok {
		return str
	}
	return ""
}

func checkParameter(dataID, group, content string) error {
	if dataID == "" || containsWhitespace(dataID) {
		return errors.New("invalid dataID")
	}
	if group == "" || containsWhitespace(group) {
		return errors.New("invalid group")
	}
	if content == "" {
		return errors.New("invalid content")
	}
	return nil
}

func generateMD5CacheKey(dataID, group string) string {
	key := group + "/" + dataID
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
