package model

import (
	"gdiamond/util/stringutil"
	"time"
)

//ConfigInfo config info model
type ConfigInfo struct {
	Group   string
	DataID  string
	Content string
	MD5     string
	//database primary key
	ID           int64
	LastModified time.Time
}

//NewConfigInfo a little bit confuse ???
func NewConfigInfo(dataID, group, content string, lastModified time.Time) *ConfigInfo {
	configInfo := &ConfigInfo{Group: group, DataID: dataID, Content: content, LastModified: lastModified}
	if content != "" {
		configInfo.MD5 = stringutil.GetMd5(content)
	}
	return configInfo
}
