/**
TODO http连接复用，参数设置细化
*/
package processor

import (
	"bufio"
	"errors"
	"fmt"
	"gdiamond/client/internal/configinfo"
	"gdiamond/client/internal/logger"
	"gdiamond/util/fileutil"
	"gdiamond/util/urlutil"
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	defaultDomainName         = "gdiamond.namesrv.net"
	dailyDomainName           = "gdiamond.namesrv.net"
	configHTTPURIFile         = "namesrv/addrs" /** 获取ServerAddress的配置uri */
	asyncAcquireIntervalInSec = 300
)

//ServerAddressProcessor name server address processor
type ServerAddressProcessor struct {
	sync.Mutex
	diamondConfigure *configinfo.Configure
	isRun            bool
}

//NewServerAddressProcessor new
func NewServerAddressProcessor(diamondConfigure *configinfo.Configure) *ServerAddressProcessor {
	p := &ServerAddressProcessor{diamondConfigure: diamondConfigure, isRun: false}
	return p
}

//Start setup processor
func (p *ServerAddressProcessor) Start() {
	p.Lock()
	defer p.Unlock()
	if p.isRun || p.diamondConfigure == nil {
		return
	}
	p.isRun = true
	if p.diamondConfigure.IsLocalFirst() {
		p.acquireServerAddressFromLocal()
	} else {
		p.synAcquireServerAddress()
		p.asynAcquireServerAddress()
	}

}

//Stop stop processor
func (p *ServerAddressProcessor) Stop() {
	p.Lock()
	if !p.isRun {
		return
	}
	p.isRun = false
	p.Unlock()
}

func (p *ServerAddressProcessor) acquireServerAddressFromLocal() error {
	if !p.isRun {
		return errors.New("ServerAddressProcessor不在运行状态，无法同步获取服务器地址列表")
	}
	acquireCount := 0
	if p.diamondConfigure.GetDomainNameList().Size() == 0 {
		p.reloadServerAddresses()
		if p.diamondConfigure.GetDomainNameList().Size() == 0 {
			if !p.acquireServerAddressOnce(acquireCount) {
				acquireCount++
				if p.acquireServerAddressOnce(acquireCount) {
					p.storeServerAddressesToLocal()
					logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from namesrv success")

				} else {
					return errors.New("当前没有可用的服务器列表")
				}
			} else {
				p.storeServerAddressesToLocal()
				logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from namesrv success")
			}
		} else {
			logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from local server address file")
		}
	}
	return nil
}

func (p *ServerAddressProcessor) synAcquireServerAddress() error {
	if !p.isRun {
		return errors.New("ServerAddressProcessor不在运行状态，无法同步获取服务器地址列表")
	}
	acquireCount := 0
	domainNameList := p.diamondConfigure.GetDomainNameList()
	if domainNameList == nil || domainNameList.Size() == 0 {
		if !p.acquireServerAddressOnce(acquireCount) {
			acquireCount++
			if p.acquireServerAddressOnce(acquireCount) {
				// 存入本地文件
				p.storeServerAddressesToLocal()
				logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from namesrv success")
			} else {
				logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from local server address file")
				p.reloadServerAddresses()
				if domainNameList.Size() == 0 {
					return errors.New("当前没有可用的服务器列表，请检查~/diamond/ServerAddress文件")
				}
			}
		} else {
			logger.Logger.WithFields(logrus.Fields{}).Debug("get server address from namesrv success")
			// 存入本地文件
			p.storeServerAddressesToLocal()
		}
	}
	return nil
}

func (p *ServerAddressProcessor) storeServerAddressesToLocal() {
	domainNameList := p.diamondConfigure.GetDomainNameList()
	filePath := generateLocalFilePath(p.diamondConfigure.GetFilePath(), "ServerAddress")
	f, err := fileutil.CreateFileIfNessary(filePath)
	defer f.Close()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err":            err.Error(),
			"domainNameList": domainNameList,
			"filePath":       filePath,
		}).Error("storeServerAddressesToLocal failed")
		return
	}
	w := bufio.NewWriter(f)
	domainNameList.Each(func(index int, value interface{}) {
		serveraddress, _ := value.(string)
		fmt.Fprintln(w, serveraddress)
	})
	w.Flush()
}

/**
 * 获取diamond服务器地址列表
 *
 * @param acquireCount
 *            根据0或1决定从日常或线上获取
 * @return
 */
func (p *ServerAddressProcessor) acquireServerAddressOnce(acquireCount int) bool {
	var configServerAddress string
	var port int
	if p.diamondConfigure.GetConfigServerAddress() != "" {
		configServerAddress = p.diamondConfigure.GetConfigServerAddress()
		port = p.diamondConfigure.GetConfigServerPort()
	} else {
		if acquireCount == 0 {
			configServerAddress = defaultDomainName
			port = configinfo.DefaultPort
		} else {
			configServerAddress = dailyDomainName
			port = configinfo.DefaultPort
		}
	}
	onceTimeOut := p.diamondConfigure.GetOnceTimeout()
	client := &http.Client{Timeout: time.Duration(onceTimeOut) * time.Millisecond}
	apiURL := urlutil.GetURL(configServerAddress, port, configHTTPURIFile)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err":            err.Error(),
			"domainNameList": apiURL,
		}).Error("NewRequest failed")
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"err": err.Error(),
			"req": req,
		}).Error("client.Do failed")
		return false
	}

	if err == nil {
		statusCode := resp.StatusCode
		if statusCode == http.StatusOK {
			newDomainNameList := dll.New()
			rd := bufio.NewReader(resp.Body)
			for {
				line, err := rd.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				} else {
					address := strings.TrimSpace(line)
					newDomainNameList.Add(address)
				}
			}
			resp.Body.Close()
			if newDomainNameList.Size() > 0 {
				logger.Logger.WithFields(logrus.Fields{
					"newDomainNameList": newDomainNameList,
				}).Debug("update DomainNameList")
				p.diamondConfigure.SetDomainNameList(newDomainNameList)
				return true
			}
		}
	}
	logger.Logger.WithFields(logrus.Fields{}).Error("no useful DomainNameList")
	return false
}

func (p *ServerAddressProcessor) asynAcquireServerAddress() {
	ticker := time.NewTicker(time.Second * time.Duration(asyncAcquireIntervalInSec))
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			if !p.isRun {
				logger.Logger.WithFields(logrus.Fields{}).Error("ServerAddressProcessor isn't running,can't get name server domain list")
				continue
			}
			acquireCount := 0
			if !p.acquireServerAddressOnce(acquireCount) {
				acquireCount++
				if p.acquireServerAddressOnce(acquireCount) {
					// 存入本地文件
					p.storeServerAddressesToLocal()
				}
			} else {
				// 存入本地文件
				p.storeServerAddressesToLocal()
			}
		}
	}()
}

func (p *ServerAddressProcessor) reloadServerAddresses() {
	filePath := generateLocalFilePath(p.diamondConfigure.GetFilePath(), "ServerAddress")
	if !fileutil.IsExist(filePath) {
		return
	}
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return
	}

	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil || io.EOF == err {
			break
		} else {
			address := strings.TrimSpace(string(line))
			p.diamondConfigure.GetDomainNameList().Add(address)
		}
	}
}

func generateLocalFilePath(directory, fileName string) string {
	if directory == "" {
		directory = os.Getenv("user.home")
	}
	return filepath.Join(directory, fileName)
}
