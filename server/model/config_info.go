package model

import "gdiamond/server/common"

type ConfigInfo struct {
	Group   string
	DataId  string
	Content string
	MD5     string
	ID      int64
}

func NewConfigInfo(dataId, group, content string) *ConfigInfo {
	configInfo := &ConfigInfo{Group: group, DataId: dataId, Content: content}
	if content != "" {
		configInfo.MD5 = common.GetMd5(content)
	}
	return configInfo
}
