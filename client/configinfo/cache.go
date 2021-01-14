package configinfo

import "sync/atomic"

type CacheData struct {
	DataId                 string
	group                  string
	md5                    string
	LastModifiedHeader     string
	DomainNamePos          int
	LocalConfigInfoFile    string
	LocalConfigInfoVersion int64
	useLocalConfigInfo     bool
	fetchCounter           int64
}

func NewCacheData(dataId, group string) *CacheData {
	return &CacheData{DataId: dataId, group: group, DomainNamePos: 0, useLocalConfigInfo: false, fetchCounter: 0}
}

func (c *CacheData) GetFetchCount() int64 {
	return atomic.LoadInt64(&c.fetchCounter)
}

func (c *CacheData) IncrementFetchCountAndGet() int64 {
	return atomic.AddInt64(&c.fetchCounter, 1)
}

func (c *CacheData) UseLocalConfigInfo() bool {
	return c.useLocalConfigInfo
}

func (c *CacheData) Group() string {
	return c.group
}

func (c *CacheData) MD5() string {
	return c.md5
}
