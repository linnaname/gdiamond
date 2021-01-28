package model

import (
	"gdiamond/server/common"
	"time"
)

type ConfigInfo struct {
	Group        string
	DataId       string
	Content      string
	MD5          string
	ID           int64
	LastModified time.Time
}

func NewConfigInfo(dataId, group, content string, lastModified time.Time) *ConfigInfo {
	configInfo := &ConfigInfo{Group: group, DataId: dataId, Content: content, LastModified: lastModified}
	if content != "" {
		configInfo.MD5 = common.GetMd5(content)
	}
	return configInfo
}
