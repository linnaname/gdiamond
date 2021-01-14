package configinfo

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WORD_SEPARATOR       = ","
	LINE_SEPARATOR       = "|"
	HTTP_URI_FILE        = "/url"
	PROBE_MODIFY_REQUEST = "Probe-Modify-Request"
)

var isRun = false
var cache sync.Map

type Subscriber struct {
	domainNamePos    int64
	diamondConfigure *Configure
}

func New() (*Subscriber, error) {
	s := &Subscriber{}
	s.diamondConfigure, _ = NewConfigure()
	return s, nil
}

func (s *Subscriber) checkDiamondServerConfigInfo() error {
	updateDataIdGroupPairs, err := s.checkUpdateDataIds(s.diamondConfigure.GetReceiveWaitTime())
	if err != nil {
		return err
	}

	if nil == updateDataIdGroupPairs || updateDataIdGroupPairs.Size() == 0 {
		log.Println("没有被修改的DataID")
		return nil
	}
	// 对于每个发生变化的DataID，都请求一次对应的配置信息

	for _, freshDataIdGroupPair := range updateDataIdGroupPairs.Values() {
		freshDataIdGroupPairStr := freshDataIdGroupPair.(string)
		middleIndex := strings.Index(freshDataIdGroupPairStr, WORD_SEPARATOR)
		if middleIndex == -1 {
			continue
		}
		freshDataId := freshDataIdGroupPairStr[0 : middleIndex-1]
		freshGroup := freshDataIdGroupPairStr[middleIndex+1:]
		value, ok := cache.Load(freshDataId)
		if !ok || value == nil {
			continue
		}
		cacheDatas, _ := value.(sync.Map)
		val, ok := cacheDatas.Load(freshGroup)
		if !ok || val == nil {
			continue
		}
		cacheData, _ := val.(CacheData)
		s.receiveConfigInfo(&cacheData)
	}
	return nil
}

/**
* 向DiamondServer请求dataId对应的配置信息，并将结果抛给客户的监听器
*
 */
func (s *Subscriber) receiveConfigInfo(cacheData *CacheData) {

}

func (s *Subscriber) checkUpdateDataIds(timeout int) (*hashset.Set, error) {
	if !isRun {
		return nil, errors.New("subscriber不在运行状态中，无法获取修改过的DataID列表")
	}

	probeUpdateString := getProbeUpdateString()
	if probeUpdateString == "" {
		return nil, errors.New("getProbeUpdateString is empty")
	}
	waitTime := 0
	for 0 == timeout || timeout > waitTime {
		onceTimeOut := s.getOnceTimeOut(waitTime, timeout)
		waitTime += onceTimeOut

		client := &http.Client{Timeout: onceTimeout * time.Millisecond}
		pos := atomic.LoadInt64(&s.domainNamePos)
		domainName, _ := s.diamondConfigure.GetDomainNameList().Get(int(pos))
		domainNamePort := getUrl(domainName.(string), s.diamondConfigure.GetPort())

		params := url.Values{}
		params.Set(PROBE_MODIFY_REQUEST, probeUpdateString)
		req, _ := http.NewRequest("POST", domainNamePort, strings.NewReader(params.Encode()))
		resp, err := client.Do(req)
		if err != nil {
			log.Println("未知异常", err)
		}
		statusCode := resp.StatusCode
		switch statusCode {
		case http.StatusOK:
			return getUpdateDataIds(resp), nil
		case http.StatusServiceUnavailable:
			s.rotateToNextDomain()
		default:
			log.Println("获取修改过的DataID列表的请求回应的HTTP State: ", statusCode)
			s.rotateToNextDomain()
		}
		resp.Body.Close()
	}

	pos := atomic.LoadInt64(&s.domainNamePos)
	domainName, _ := s.diamondConfigure.GetDomainNameList().Get(int(pos))
	return nil, errors.New(fmt.Sprintf("获取修改过的DataID列表超时:%v, 超时时间为:%v", domainName, timeout))
}

func getProbeUpdateString() string {
	probeModifyBuilder := strings.Builder{}
	cache.Range(func(key, value interface{}) bool {
		dataId, _ := key.(string)
		if value == nil {
			return true
		}
		groupCache, ok := value.(sync.Map)
		if !ok {
			return true
		}
		groupCache.Range(func(key, value interface{}) bool {
			if value == nil {
				return true
			}
			data, ok := value.(CacheData)
			// 非使用本地配置，才去diamond server检查
			if !ok {
				return true
			}
			if !data.UseLocalConfigInfo() {
				probeModifyBuilder.WriteString(dataId)
				probeModifyBuilder.WriteString(WORD_SEPARATOR)

				if data.Group() != "" {
					probeModifyBuilder.WriteString(data.Group())
					probeModifyBuilder.WriteString(WORD_SEPARATOR)
				} else {
					probeModifyBuilder.WriteString(WORD_SEPARATOR)
				}

				if data.MD5() != "" {
					probeModifyBuilder.WriteString(data.MD5())
					probeModifyBuilder.WriteString(LINE_SEPARATOR)
				} else {
					probeModifyBuilder.WriteString(LINE_SEPARATOR)
				}
			}
			return true
		})
		return true
	})
	return probeModifyBuilder.String()
}

/**
waitTime 本次查询已经耗费的时间(已经查询的多次HTTP耗费的时间)
timeout 本次查询总的可耗费时间(可供多次HTTP查询使用)
return 本次HTTP查询能够使用的时间
*/
func (s *Subscriber) getOnceTimeOut(waitTime, timeout int) int {
	onceTimeOut := s.diamondConfigure.GetOnceTimeout()
	remainTime := timeout - waitTime
	if onceTimeOut > remainTime {
		onceTimeOut = remainTime
	}
	return onceTimeOut
}

func getUrl(domainName string, port int) string {
	apiUrl := domainName + ":" + strconv.Itoa(port)
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = HTTP_URI_FILE
	urlStr := u.String()
	return urlStr
}

func (s *Subscriber) rotateToNextDomain() {

}

func getUpdateDataIds(resp *http.Response) *hashset.Set {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return convertStringToSet(string(body))
}

func convertStringToSet(modifiedDataIdsString string) *hashset.Set {
	if modifiedDataIdsString == "" {
		return nil
	}
	modifiedDataIdsString, err := url.QueryUnescape(modifiedDataIdsString)
	if err != nil {
		log.Println("解码modifiedDataIdsString出错", err)
	}

	if modifiedDataIdsString != "" {
		if strings.HasPrefix(modifiedDataIdsString, "OK") {
			log.Println("探测的返回结果:" + modifiedDataIdsString)
		} else {
			log.Println("探测到数据变化:" + modifiedDataIdsString)
		}
	}

	modifiedDataIdSet := hashset.New()
	modifiedDataIdStrings := strings.Split(modifiedDataIdsString, LINE_SEPARATOR)
	for _, modifiedDataIdString := range modifiedDataIdStrings {
		if modifiedDataIdString != "" {
			modifiedDataIdSet.Add(modifiedDataIdString)
		}
	}
	return modifiedDataIdSet
}
