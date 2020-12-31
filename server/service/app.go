package service

import (
	"crypto/md5"
	"fmt"
	"gdiamond/server/model"
	"sync"
)

var contentMD5Cache sync.Map

func UpdateMD5Cache(configInfo *model.ConfigInfo) {
	key := generateMD5CacheKey(configInfo.DataId, configInfo.Group)
	md5 := fmt.Sprintf("%x", md5.Sum([]byte(configInfo.Content)))
	contentMD5Cache.Store(key, md5)
}

func generateMD5CacheKey(dataId, group string) string {
	key := group + "/" + dataId
	return key
}
