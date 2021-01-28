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

func FindConfigInfo(dataId, group string) (*model.ConfigInfo, error) {
	return findConfigInfo(dataId, group)
}

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
		return findConfigInfoByDataId(pageNo, pageSize, dataId)
	} else if dataId != "" && group == "" {
		return findConfigInfoByGroup(pageNo, pageSize, group)
	} else {
		return findAllConfigInfo(pageNo, pageSize)
	}
}

func FindConfigInfoLike(pageNo, pageSize int, dataId, group string) (*model.Page, error) {
	return FindConfigInfoLike(pageNo, pageSize, dataId, group)
}

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
	value := i2Str(configInfo.(model.ConfigInfo).MD5)
	return value
}

func GetLastModified(dataId, group string) string {
	locker.Lock()
	defer locker.Unlock()
	key := generateMD5CacheKey(dataId, group)
	value, loaded := cache.Load(key)
	if value == nil || !loaded {
		return ""
	}
	configInfo, _ := value.(model.ConfigInfo)
	lastModified := configInfo.LastModified
	return lastModified.String()
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

func GetConfigInfoPath(dataId, group string) string {
	builder := strings.Builder{}
	builder.WriteString("/")
	builder.WriteString(BASE_DIR)
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

func UpdateMD5Cache(configInfo *model.ConfigInfo) {
	key := generateMD5CacheKey(configInfo.DataId, configInfo.Group)
	md5 := common.GetMd5(configInfo.Content)
	configInfo.MD5 = md5
	cache.Store(key, configInfo)
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
