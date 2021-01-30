package configinfo

import (
	"gdiamond/util/fileutil"
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
)

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
	DEFAULT_PORT          = 8080
	MaxUint               = ^uint(0)
	POLLING_INTERVAL_TIME = 15 // 秒
	DEFAULT_GROUP         = "DEFAULT_GROUP"
)

func NewConfigure() (*Configure, error) {
	filePath := fileutil.GetCurrentDirectory() + "/gdiamond"
	err := fileutil.CreateDirIfNessary(filePath)
	return &Configure{filePath: filePath, domainNameList: dll.New(), localFirst: false, configServerPort: DEFAULT_PORT, port: 1210,
		onceTimeout: 2000, receiveWaitTime: 2000 * 5, retrieveDataRetryTimes: int(MaxUint>>1) / 10}, err
}

func (c *Configure) GetFilePath() string {
	return c.filePath
}

/**
 * 获取对于一个DiamondServer所对应的查询一个DataID对应的配置信息的Timeout时间<br>
 * 即一次HTTP请求的超时时间<br>
 * 单位：毫秒<br>
 */
func (c *Configure) GetOnceTimeout() int {
	return c.onceTimeout
}

func (c *Configure) GetPort() int {
	return c.port
}

/**
 * 同步查询一个DataID的最长等待时间<br>
 * 实际最长等待时间小于receiveWaitTime + min(connectionTimeout, onceTimeout)
 *
 * @return
 */
func (c *Configure) GetReceiveWaitTime() int {
	return c.receiveWaitTime
}

func (c *Configure) GetRetrieveDataRetryTimes() int {
	return c.retrieveDataRetryTimes
}

/**
 * 获取轮询的间隔时间。单位：秒<br>
 * 此间隔时间代表轮询查找一次配置信息的间隔时间，对于容灾相关，请设置短一些；<br>
 * 对于其他不可变的配置信息，请设置长一些
 */
func (c *Configure) GetPollingIntervalTime() int64 {
	return c.pollingIntervalTime
}

/**
 * 获取当前支持的所有的DiamondServer域名列表
 */
func (c *Configure) GetDomainNameList() *dll.List {
	return c.domainNameList
}

/**
 * 设置当前支持的所有的DiamondServer域名列表，当设置了域名列表后，缺省的域名列表将失效
 */
func (c *Configure) SetDomainNameList(domainNameList *dll.List) {
	if nil == domainNameList {
		return
	}
	c.domainNameList = domainNameList
}

/**
 * 添加一个DiamondServer域名，当设置了域名列表后，缺省的域名列表将失效
 */
func (c *Configure) AddDomainName(domainName string) {
	if "" == domainName {
		return
	}
	c.domainNameList.Add(domainName)
}

/**
 * 设置轮询的间隔时间。单位：秒
 */
func (c *Configure) SetPollingIntervalTime(pollingIntervalTime int64) {
	if pollingIntervalTime < POLLING_INTERVAL_TIME {
		return
	}
	c.pollingIntervalTime = pollingIntervalTime
}

func (c *Configure) IsLocalFirst() bool {
	return c.localFirst
}

func (c *Configure) GetConfigServerAddress() string {
	return c.configServerAddress
}

func (c *Configure) GetConfigServerPort() int {
	return c.configServerPort
}
