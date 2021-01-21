/**
TODO http连接复用，参数设置细化
TODO serveraddres请求namesrv
*/
package processor

import (
	"bufio"
	"errors"
	"fmt"
	"gdiamond/client/configinfo"
	"gdiamond/util/fileutil"
	"gdiamond/util/urlutil"
	dll "github.com/emirpasic/gods/lists/doublylinkedlist"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_DOMAINNAME       = "http://127.0.0.1"
	DAILY_DOMAINNAME         = "http://127.0.0.1"
	CONFIG_HTTP_URI_FILE     = "url" /** 获取ServerAddress的配置uri */
	asynAcquireIntervalInSec = 300
)

type ServerAddressProcessor struct {
	sync.Mutex
	diamondConfigure *configinfo.Configure
	isRun            bool
}

func NewServerAddressProcessor(diamondConfigure *configinfo.Configure) *ServerAddressProcessor {
	p := &ServerAddressProcessor{diamondConfigure: diamondConfigure, isRun: false}
	return p
}

func (p *ServerAddressProcessor) Start() {
	p.Lock()
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
	p.Unlock()
}

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
					log.Println("在同步获取服务器列表时，向日常ConfigServer服务器获取到了服务器列表")
				} else {
					return errors.New("当前没有可用的服务器列表")
				}
			} else {
				p.storeServerAddressesToLocal()
				log.Println("在同步获取服务器列表时，向线上ConfigServer服务器获取到了服务器列表")
			}
		} else {
			log.Println("在同步获取服务器列表时，由于本地指定了服务器列表，不向ConfigServer服务器同步获取服务器列表")
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
				log.Println("在同步获取服务器列表时，向日常ConfigServer服务器获取到了服务器列表")
			} else {
				log.Println("从本地获取Diamond地址列表")
				p.reloadServerAddresses()
				if domainNameList.Size() == 0 {
					return errors.New("当前没有可用的服务器列表，请检查~/diamond/ServerAddress文件")
				}
			}
		} else {
			log.Println("在同步获取服务器列表时，向线上ConfigServer服务器获取到了服务器列表")
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
		log.Println("存储服务器地址到本地文件失败", err)
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
			configServerAddress = DEFAULT_DOMAINNAME
			port = configinfo.DEFAULT_PORT
		} else {
			configServerAddress = DAILY_DOMAINNAME
			port = configinfo.DEFAULT_PORT
		}
	}
	onceTimeOut := p.diamondConfigure.GetOnceTimeout()
	client := &http.Client{Timeout: time.Duration(onceTimeOut) * time.Millisecond}
	apiUrl := urlutil.GetUrl(configServerAddress, port, CONFIG_HTTP_URI_FILE)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println("NewRequest error", err)
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("没有可用的新服务器列表", err)
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
				log.Println("更新使用的服务器列表")
				p.diamondConfigure.SetDomainNameList(newDomainNameList)
				return true
			}
		}
	}
	log.Println("没有可用的新服务器列表", err)
	return false
}

func (p *ServerAddressProcessor) asynAcquireServerAddress() {
	ticker := time.NewTicker(time.Second * time.Duration(asynAcquireIntervalInSec))
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			if !p.isRun {
				log.Println("ServerAddressProcessor不在运行状态，无法异步获取服务器地址列表")
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
			//FIXME 已经是定时任务了有必要这里多递归吗？
			//p.asynAcquireServerAddress()
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
