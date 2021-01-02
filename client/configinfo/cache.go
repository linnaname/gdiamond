package configinfo

import "sync/atomic"

type CacheData struct {
	DataId                 string
	Group                  string
	MD5                    string
	LastModifiedHeader     string
	DomainNamePos          int
	LocalConfigInfoFile    string
	LocalConfigInfoVersion int64
	UseLocalConfigInfo     bool
	fetchCounter           int64
}

func NewCacheData(dataId, group string) *CacheData {
	return &CacheData{DataId: dataId, Group: group, DomainNamePos: 0, UseLocalConfigInfo: false, fetchCounter: 0}
}

func (c *CacheData) GetFetchCount() int64 {
	return atomic.LoadInt64(&c.fetchCounter)
}

func (c *CacheData) IncrementFetchCountAndGet() int64 {
	return atomic.AddInt64(&c.fetchCounter, 1)
}
