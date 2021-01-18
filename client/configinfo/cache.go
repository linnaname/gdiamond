package configinfo

import "sync/atomic"

type CacheData struct {
	dataId                 string
	group                  string
	md5                    string
	lastModifiedHeader     string
	domainNamePos          int
	localConfigInfoFile    string
	localConfigInfoVersion int64
	useLocalConfigInfo     bool
	//couter of success fetch
	fetchCounter int64
}

func NewCacheData(dataId, group string) *CacheData {
	return &CacheData{dataId: dataId, group: group, domainNamePos: 0, useLocalConfigInfo: false, fetchCounter: 0}
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

func (c *CacheData) SetUseLocalConfigInfo(useLocalConfigInfo bool) {
	c.useLocalConfigInfo = useLocalConfigInfo
}

func (c *CacheData) Group() string {
	return c.group
}

func (c *CacheData) DataId() string {
	return c.dataId
}

func (c *CacheData) MD5() string {
	return c.md5
}

func (c *CacheData) SetMD5(md5 string) {
	c.md5 = md5
}

func (c *CacheData) GetLastModifiedHeader() string {
	return c.lastModifiedHeader
}

func (c *CacheData) SetLastModifiedHeader(lastModifiedHeader string) {
	c.lastModifiedHeader = lastModifiedHeader
}

func (c *CacheData) GetLocalConfigInfoFile() string {
	return c.localConfigInfoFile
}

func (c *CacheData) SetLocalConfigInfoFile(localConfigInfoFile string) {
	c.localConfigInfoFile = localConfigInfoFile
}

func (c *CacheData) GetLocalConfigInfoVersion() int64 {
	return c.localConfigInfoVersion
}

func (c *CacheData) SetLocalConfigInfoVersion(localConfigInfoVersion int64) {
	c.localConfigInfoVersion = localConfigInfoVersion
}
