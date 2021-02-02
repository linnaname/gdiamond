package configinfo

import "sync/atomic"

//CacheData  only in memory
type CacheData struct {
	dataId                 string
	group                  string
	md5                    string
	content                string
	lastModifiedHeader     string
	domainNamePos          int
	localConfigInfoFile    string
	localConfigInfoVersion int64
	useLocalConfigInfo     bool
	//counter of success fetch
	fetchCounter int64
}

//NewCacheData new
func NewCacheData(dataId, group string) *CacheData {
	return &CacheData{dataId: dataId, group: group, domainNamePos: 0, useLocalConfigInfo: false, fetchCounter: 0}
}

//GetFetchCount fetCounter getter
func (c *CacheData) GetFetchCount() int64 {
	return atomic.LoadInt64(&c.fetchCounter)
}

//IncrementFetchCountAndGet atomic fetchCounter increment and getter
func (c *CacheData) IncrementFetchCountAndGet() int64 {
	return atomic.AddInt64(&c.fetchCounter, 1)
}

//UseLocalConfigInfo useLocalConfigInfo getter
func (c *CacheData) UseLocalConfigInfo() bool {
	return c.useLocalConfigInfo
}

//SetUseLocalConfigInfo useLocalConfigInfo setter
func (c *CacheData) SetUseLocalConfigInfo(useLocalConfigInfo bool) {
	c.useLocalConfigInfo = useLocalConfigInfo
}

//Group group getter
func (c *CacheData) Group() string {
	return c.group
}

//DataId md5 getter
func (c *CacheData) DataId() string {
	return c.dataId
}

//MD5 md5 getter
func (c *CacheData) MD5() string {
	return c.md5
}

//SetMD5 md5 setter
func (c *CacheData) SetMD5(md5 string) {
	c.md5 = md5
}

//GetLastModifiedHeader lastModifiedHeader getter
func (c *CacheData) GetLastModifiedHeader() string {
	return c.lastModifiedHeader
}

//SetLastModifiedHeader lastModifiedHeader setter
func (c *CacheData) SetLastModifiedHeader(lastModifiedHeader string) {
	c.lastModifiedHeader = lastModifiedHeader
}

//GetLocalConfigInfoFile localConfigInfoFile getter
func (c *CacheData) GetLocalConfigInfoFile() string {
	return c.localConfigInfoFile
}

//SetLocalConfigInfoFile localConfigInfoFile setter
func (c *CacheData) SetLocalConfigInfoFile(localConfigInfoFile string) {
	c.localConfigInfoFile = localConfigInfoFile
}

//GetLocalConfigInfoVersion localConfigInfoVersion getter
func (c *CacheData) GetLocalConfigInfoVersion() int64 {
	return c.localConfigInfoVersion
}

//SetLocalConfigInfoVersion localConfigInfoVersion setter
func (c *CacheData) SetLocalConfigInfoVersion(localConfigInfoVersion int64) {
	c.localConfigInfoVersion = localConfigInfoVersion
}
