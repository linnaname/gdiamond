package configinfo

import "os"
import dll "github.com/emirpasic/gods/lists/doublylinkedlist"

type Configure struct {
	filePath       string // 本地数据保存路径
	domainNameList *dll.List
}

const (
	onceTimeout         = 2000 //毫秒，获取对于一个DiamondServer所对应的查询一个DataID对应的配置信息的Timeout时间
	pollingIntervalTime = 15   //异步查询的间隔时间
	port                = 8080
	receiveWaitTime     = onceTimeout * 5 // 同步查询一个DataID所花费的时间
)

func NewConfigure() (*Configure, error) {
	filePath := os.Getenv("user.home") + "/diamond"
	err := createDirIfNessary(filePath)
	return &Configure{filePath: filePath}, err
}

func createDirIfNessary(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

/**
 * 获取对于一个DiamondServer所对应的查询一个DataID对应的配置信息的Timeout时间<br>
 * 即一次HTTP请求的超时时间<br>
 * 单位：毫秒<br>
 */
func (c *Configure) GetOnceTimeout() int {
	return onceTimeout
}

func (c *Configure) GetPort() int {
	return port
}

/**
 * 同步查询一个DataID的最长等待时间<br>
 * 实际最长等待时间小于receiveWaitTime + min(connectionTimeout, onceTimeout)
 *
 * @return
 */
func (c *Configure) GetReceiveWaitTime() int {
	return receiveWaitTime
}

/**
 * 获取当前支持的所有的DiamondServer域名列表
 *
 * @return
 */
func (c *Configure) GetDomainNameList() *dll.List {
	return c.domainNameList
}

/**
 * 设置当前支持的所有的DiamondServer域名列表，当设置了域名列表后，缺省的域名列表将失效
 *
 * @param domainNameList
 */
func (c *Configure) SetDomainNameList(domainNameList *dll.List) {
	if nil == domainNameList {
		return
	}
	c.domainNameList = dll.New()
}

/**
 * 添加一个DiamondServer域名，当设置了域名列表后，缺省的域名列表将失效
 *
 * @param domainName
 */
func (c *Configure) AddDomainName(domainName string) {
	if "" == domainName {
		return
	}
	c.domainNameList.Add(domainName)
}
