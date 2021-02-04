package configinfo

import (
	"gdiamond/util/fileutil"
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

//Configure client configure
type Configure struct {
	filePath               string // 本地数据保存路径
	domainNameList         *dll.List
	pollingIntervalTime    int64
	localFirst             bool
	configServerAddress    string
	configServerPort       int
	port                   int
	onceTimeout            int //获取对于一个DiamondServer所对应的查询一个DataID对应的配置信息的Timeout时间(毫秒)
	receiveWaitTime        int // 同步查询一个DataID所花费的时间
	retrieveDataRetryTimes int //获取数据时的重试次数
}

const (
	//DefaultPort default http server post
	DefaultPort = 8080
	//MaxUint max num of unit type
	MaxUint = ^uint(0)
	//PollingIntervalTime seconds
	PollingIntervalTime = 15
	//DefaultGroup default group
	DefaultGroup = "DEFAULT_GROUP"
)

//NewConfigure filePath need to modified,but not yet
func NewConfigure() (*Configure, error) {
	filePath := fileutil.GetCurrentDirectory() + "/gdiamond"
	err := fileutil.CreateDirIfNecessary(filePath)
	return &Configure{filePath: filePath, domainNameList: dll.New(), localFirst: false, configServerPort: DefaultPort, port: 1210,
		onceTimeout: 2000, receiveWaitTime: 2000 * 5, retrieveDataRetryTimes: int(MaxUint>>1) / 10}, err
}

//GetFilePath filePath getter,where client store local file
func (c *Configure) GetFilePath() string {
	return c.filePath
}

//GetOnceTimeout 获取对于一个DiamondServer所对应的查询一个DataID对应的配置信息的Timeout时间,即一次HTTP请求的超时时间 单位：毫秒
func (c *Configure) GetOnceTimeout() int {
	return c.onceTimeout
}

//GetPort port getter
func (c *Configure) GetPort() int {
	return c.port
}

//GetReceiveWaitTime 同步查询一个DataID的最长等待时间,实际最长等待时间小于receiveWaitTime + min(connectionTimeout, onceTimeout)
func (c *Configure) GetReceiveWaitTime() int {
	return c.receiveWaitTime
}

//GetRetrieveDataRetryTimes 重试次数
func (c *Configure) GetRetrieveDataRetryTimes() int {
	return c.retrieveDataRetryTimes
}

//GetPollingIntervalTime 获取轮询的间隔时间。单位：秒,此间隔时间代表轮询查找一次配置信息的间隔时间，对于容灾相关，请设置短一些；<br>
//对于其他不可变的配置信息，请设置长一些
func (c *Configure) GetPollingIntervalTime() int64 {
	return c.pollingIntervalTime
}

//GetDomainNameList 获取当前支持的所有的DiamondServer域名列表
func (c *Configure) GetDomainNameList() *dll.List {
	return c.domainNameList
}

//SetDomainNameList 设置当前支持的所有的DiamondServer域名列表，当设置了域名列表后，缺省的域名列表将失效
func (c *Configure) SetDomainNameList(domainNameList *dll.List) {
	if nil == domainNameList {
		return
	}
	c.domainNameList = domainNameList
}

//AddDomainName 添加一个DiamondServer域名，当设置了域名列表后，缺省的域名列表将失效
func (c *Configure) AddDomainName(domainName string) {
	if "" == domainName {
		return
	}
	c.domainNameList.Add(domainName)
}

//SetPollingIntervalTime 设置轮询的间隔时间。单位：秒
func (c *Configure) SetPollingIntervalTime(pollingIntervalTime int64) {
	if pollingIntervalTime < PollingIntervalTime {
		return
	}
	c.pollingIntervalTime = pollingIntervalTime
}

//IsLocalFirst localFirst getter,优先使用local file
func (c *Configure) IsLocalFirst() bool {
	return c.localFirst
}

//GetConfigServerAddress configServerAddress getter
func (c *Configure) GetConfigServerAddress() string {
	return c.configServerAddress
}

//GetConfigServerPort configServerPort getter
func (c *Configure) GetConfigServerPort() int {
	return c.configServerPort
}
