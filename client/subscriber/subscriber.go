/**
TODO http连接复用，精细参数设置
TODO 优雅关闭procc等
*/

package subscriber

import (
	"compress/gzip"
	"errors"
	"fmt"
	"gdiamond/client/internal/configinfo"
	"gdiamond/client/internal/processor"
	"gdiamond/client/internal/simplecache"
	"gdiamond/client/listener"
	"gdiamond/util/maputil"
	"gdiamond/util/stringutil"
	"gdiamond/util/urlutil"
	"github.com/emirpasic/gods/sets/hashset"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	wordSeparator      = ","
	lineSeparator      = ";"
	probeModifyRequest = "Probe-Modify-Request"
	dataID             = "dataId"
	GROUP              = "group"
	CONTENT            = "content"
	IfModifiedSince    = "If-Modified-Since"
	ContentMd5         = "Content-MD5"
	AcceptEncoding     = "Accept-Encoding"
	ContentEncoding    = "Content-Encoding"
	LastModified       = "Last-Modified"
	SpacingInterval    = "client-spacing-interval"
	DataDir            = "data"     // local dir need to watch
	SnapshotDir        = "snapshot" // last time succeed snapshot  dir
	GetConfigUrl       = "diamond-server/config"
	PublishConfigUrl   = "diamond-server/publishConfig"
	GetProbeModifyUrl  = "diamond-server/getProbeModify"
)

//Subscriber client subscriber
type Subscriber struct {
	//lock to keep start and stop goroutine-safe
	sync.Mutex
	//random  domainName index
	domainNamePos    int64
	diamondConfigure *configinfo.Configure
	//first time check local file
	bFirstCheck  bool
	contentCache *simplecache.SimpleCache
	cache        sync.Map
	//subscriber is started
	isRun bool

	//local file processor
	localConfigInfoProcessor *processor.LocalConfigInfoProcessor
	//server address processor
	serverAddressProcessor *processor.ServerAddressProcessor
	//last time succeed snapshot  processor
	snapshotConfigInfoProcessor *processor.SnapshotConfigInfoProcessor

	subscriberListener listener.SubscriberListener
}

func newSubscriber(subscriberListener listener.SubscriberListener) (*Subscriber, error) {
	s := &Subscriber{bFirstCheck: true, isRun: false}
	s.diamondConfigure, _ = configinfo.NewConfigure()
	s.contentCache = simplecache.New(s.diamondConfigure.GetPollingIntervalTime())
	s.subscriberListener = subscriberListener
	return s, nil
}

//Start start subscriber
//start local file/last time succeed snapshot/server address processor
//random domain name index
func (s *Subscriber) Start() {
	s.Lock()
	defer s.Unlock()
	if s.isRun {
		return
	}
	s.localConfigInfoProcessor = processor.NewLocalConfigInfoProcessor()
	s.localConfigInfoProcessor.Start(s.diamondConfigure.GetFilePath() + "/" + DataDir)

	s.serverAddressProcessor = processor.NewServerAddressProcessor(s.diamondConfigure)
	s.serverAddressProcessor.Start()

	s.snapshotConfigInfoProcessor = processor.NewSnapshotConfigInfoProcessor(s.diamondConfigure.GetFilePath() + "/" + SnapshotDir)

	s.randomDomainNamePos()
	s.isRun = true
	s.diamondConfigure.SetPollingIntervalTime(configinfo.PollingIntervalTime)
	s.rotateCheckConfigInfo()

}

//Close close subscriber
//1.close local file/last time succeed snapshot/server address processor
//2.clear cache
//Is there a better way to close gracefully???
func (s *Subscriber) Close() {
	s.Lock()
	if !s.isRun {
		return
	}
	log.Println("start close subscriber")

	s.localConfigInfoProcessor.Stop()
	s.serverAddressProcessor.Stop()
	s.isRun = false
	maputil.ClearSyncMap(s.cache)
	log.Println("subscriber closed")
	s.Unlock()
}

//GetSubscriberListener getter
func (s *Subscriber) GetSubscriberListener() listener.SubscriberListener {
	return s.subscriberListener
}

//SetSubscriberListener setter
func (s *Subscriber) SetSubscriberListener(l listener.SubscriberListener) {
	s.subscriberListener = l
}

//AddDataId add dataId to watch
func (s *Subscriber) AddDataId(dataId, group string) {
	if group == "" {
		group = configinfo.DefaultGroup
	}
	value, ok := s.cache.Load(dataId)
	var cacheDatas sync.Map
	if !ok || value == nil {
		var newCacheDatas sync.Map
		actual, loaded := s.cache.LoadOrStore(dataId, newCacheDatas)
		if nil != actual && loaded {
			oldCacheDatas, _ := actual.(sync.Map)
			cacheDatas = oldCacheDatas
		} else {
			cacheDatas = newCacheDatas
		}
	} else {
		cacheDatas = value.(sync.Map)
	}
	cacheData, ok := cacheDatas.Load(group)
	if nil == cacheData || !ok {
		//s.Start()
		c := configinfo.NewCacheData(dataId, group)
		content := s.loadCacheContentFromDiskLocal(c)
		md5 := stringutil.GetMd5(content)
		c.SetMD5(md5)
		cacheDatas.LoadOrStore(group, c)
		s.cache.Store(dataId, cacheDatas)
	}
}

//RemoveDataId delete dataId to watch
func (s *Subscriber) RemoveDataId(dataId, group string) {
	if group == "" {
		group = configinfo.DefaultGroup
	}
	value, ok := s.cache.Load(dataId)
	if !ok || value == nil {
		return
	}
	cacheDatas, _ := value.(sync.Map)
	cacheDatas.Delete(group)

	log.Println("删除了DataID[" + dataId + "]中的Group: " + group)
	if maputil.LengthOfSyncMap(cacheDatas) == 0 {
		s.cache.Delete(dataId)
		log.Println("删除了DataID[" + dataId + "]")
	}
}

//GetDataIds get dataIds on watching
func (s *Subscriber) GetDataIds() *hashset.Set {
	keys := hashset.New()
	s.cache.Range(func(k, v interface{}) bool {
		keys.Add(k)
		return true
	})
	return keys
}

func (s *Subscriber) ClearCache() {
	maputil.ClearSyncMap(s.cache)
}

//GetConfigureInformation implement method
func (s *Subscriber) GetConfigureInformation(dataId, group string, timeout int) string {
	if group == "" {
		group = configinfo.DefaultGroup
	}
	cacheData := s.getCacheData(dataId, group)
	// 优先使用本地配置
	localConfig, err := s.localConfigInfoProcessor.GetLocalConfigureInformation(cacheData, true)
	if localConfig != "" {
		cacheData.IncrementFetchCountAndGet()
		s.saveSnapshot(dataId, group, localConfig)
		return localConfig
	}

	if err != nil {
		log.Println("获取本地配置文件出错", err)
	}
	result, err := s.getConfigureInfomation(dataId, group, timeout, false)
	if err == nil && result != "" {
		s.saveSnapshot(dataId, group, result)
		cacheData.IncrementFetchCountAndGet()
	}
	return result
}

//GetAvailableConfigureInformation  implement method
func (s *Subscriber) GetAvailableConfigureInformation(dataId, group string, timeout int) string {
	// 尝试先从本地和网络获取配置信息
	result := s.GetConfigureInformation(dataId, group, timeout)
	if result != "" {
		return result
	}
	return s.getSnapshotConfiginfomation(dataId, group)
}

//GetAvailableConfigureInformationFromSnapshot implement method
func (s *Subscriber) GetAvailableConfigureInformationFromSnapshot(dataId, group string, timeout int) string {
	result := s.getSnapshotConfiginfomation(dataId, group)
	if result != "" {
		return result
	}
	return s.GetConfigureInformation(dataId, group, timeout)
}

//PublishConfigureInformation implement method
func (s *Subscriber) PublishConfigureInformation(dataId, group, content string) error {
	if !s.isRun {
		return errors.New("subscriber isn't running can't publish config info")
	}
	if content == "" {
		return errors.New("content can't be empty")
	}
	timeout := s.diamondConfigure.GetReceiveWaitTime()
	waitTime := 0
	retryTimes := s.diamondConfigure.GetRetrieveDataRetryTimes()
	log.Println("设定的获取配置数据的重试次数为：", retryTimes)
	// 已经尝试过的次数
	tryCount := 0
	for 0 == timeout || timeout > waitTime {
		// 尝试次数加1
		tryCount++
		if tryCount > retryTimes+1 {
			log.Println("已经到达了设定的重试次数")
			break
		}
		log.Println(fmt.Sprintf("获取配置数据，第%v,次尝试, waitTime:%v", tryCount, waitTime))

		// 设置超时时间
		onceTimeOut := s.getOnceTimeOut(waitTime, timeout)
		waitTime += onceTimeOut

		client := &http.Client{Timeout: time.Duration(onceTimeOut) * time.Millisecond}
		pos := atomic.LoadInt64(&s.domainNamePos)
		domainNameList := s.diamondConfigure.GetDomainNameList()
		domainName, ok := domainNameList.Get(int(pos))
		if !ok || domainName == nil {
			log.Println("无可用domainName")
			s.rotateToNextDomain()
			continue
		}
		domainNamePort := urlutil.GetURL(domainName.(string), s.diamondConfigure.GetPort(), PublishConfigUrl)
		params := url.Values{}
		params.Set(dataID, dataId)
		params.Set(GROUP, group)
		params.Set(CONTENT, content)
		req, err := http.NewRequest("POST", domainNamePort, ioutil.NopCloser(strings.NewReader(params.Encode())))
		if err != nil {
			log.Println("Can't NewRequest", err)
			continue
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("未知异常", err)
			s.rotateToNextDomain()
			continue
		}
		statusCode := resp.StatusCode
		if statusCode == http.StatusOK {
			return nil
		}
		log.Println("publish config info timeoutHTTP State: ", statusCode)
		s.rotateToNextDomain()
		resp.Body.Close()
	}

	return fmt.Errorf("publish config info timeout:%v", timeout)
}

func (s *Subscriber) loadCacheContentFromDiskLocal(cacheData *configinfo.CacheData) string {
	content, _ := s.localConfigInfoProcessor.GetLocalConfigureInformation(cacheData, true)
	if content != "" {
		return content
	}
	c, _ := s.snapshotConfigInfoProcessor.GetConfigInfomation(cacheData.DataId(), cacheData.Group())
	return c
}

/**
 * 循环探测配置信息是否变化，如果变化，则再次向DiamondServer请求获取对应的配置信息
 */
func (s *Subscriber) rotateCheckConfigInfo() {
	duration := 30
	if !s.bFirstCheck {
		duration = int(s.diamondConfigure.GetPollingIntervalTime())
	}
	ticker := time.NewTicker(time.Second * time.Duration(duration))
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			if !s.isRun {
				log.Println("DiamondSubscriber不在运行状态中，退出查询循环")
				return
			}
			s.checkLocalConfigInfo()
			err := s.checkDiamondServerConfigInfo()
			if err != nil {
				log.Println("循环探测发生异常", err)
			}
			s.checkSnapshot()
		}
	}()
	s.bFirstCheck = false
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
		middleIndex := strings.Index(freshDataIdGroupPairStr, wordSeparator)
		if middleIndex == -1 {
			continue
		}
		freshDataId := freshDataIdGroupPairStr[0:middleIndex]
		freshGroup := freshDataIdGroupPairStr[middleIndex+1:]
		value, ok := s.cache.Load(freshDataId)
		if !ok || value == nil {
			continue
		}
		cacheDatas, _ := value.(sync.Map)
		val, ok := cacheDatas.Load(freshGroup)
		if !ok || val == nil {
			continue
		}
		cacheData, _ := val.(*configinfo.CacheData)
		s.receiveConfigInfo(cacheData)
	}
	return nil
}

/**
* 向DiamondServer请求dataId对应的配置信息，并将结果抛给客户的监听器
*
 */
func (s *Subscriber) receiveConfigInfo(cacheData *configinfo.CacheData) {
	go func() {
		if !s.isRun {
			log.Println("DiamondSubscriber不在运行状态中，退出查询循环")
			return
		}
		configInfo, err := s.getConfigureInfomation(cacheData.DataId(), cacheData.Group(), s.diamondConfigure.GetReceiveWaitTime(), true)
		if err != nil {
			log.Println("向Diamond服务器索要配置信息的过程抛异常", err)
			return
		}

		if configInfo == "" {
			return
		}

		if nil == s.subscriberListener {
			log.Println("null == subscriberListener")
			return
		}

		s.popConfigInfo(cacheData, configInfo)
	}()
}

/**
是否使用本地的内容cache。主动get时会使用，check触发的异步get不使用本地cache。
*/
func (s *Subscriber) getConfigureInfomation(dataId, group string, timeout int, skipContentCache bool) (string, error) {
	s.Start()
	if !s.isRun {
		return "", errors.New("DiamondSubscriber不在运行状态中，无法获取ConfigureInfomation")
	}
	if group == "" {
		group = configinfo.DefaultGroup
	}

	/**
	 * 使用带有TTL的cache，
	 */
	if !skipContentCache {
		key := makeCacheKey(dataId, group)
		content := s.contentCache.Get(key)
		if content != nil {
			return content.(string), nil
		}
	}

	waitTime := 0
	cacheData := s.getCacheData(dataId, group)
	retryTimes := s.diamondConfigure.GetRetrieveDataRetryTimes()
	log.Println("设定的获取配置数据的重试次数为：", retryTimes)
	// 已经尝试过的次数
	tryCount := 0
	for 0 == timeout || timeout > waitTime {
		// 尝试次数加1
		tryCount++
		if tryCount > retryTimes+1 {
			log.Println("已经到达了设定的重试次数")
			break
		}
		log.Println(fmt.Sprintf("获取配置数据，第%v,次尝试, waitTime:%v", tryCount, waitTime))

		// 设置超时时间
		onceTimeOut := s.getOnceTimeOut(waitTime, timeout)
		waitTime += onceTimeOut

		client := &http.Client{Timeout: time.Duration(onceTimeOut) * time.Millisecond}
		pos := atomic.LoadInt64(&s.domainNamePos)
		domainNameList := s.diamondConfigure.GetDomainNameList()
		domainName, ok := domainNameList.Get(int(pos))
		if !ok || domainName == nil {
			log.Println("无可用domainName")
			s.rotateToNextDomain()
			continue
		}
		domainNamePort := urlutil.GetURL(domainName.(string), s.diamondConfigure.GetPort(), GetConfigUrl)
		//params := url.Values{}
		//params.Set(DATAID, dataId)
		//params.Set(GROUP, group)
		req, err := http.NewRequest("GET", domainNamePort, nil)
		if err != nil {
			log.Println("Can't NewRequest", err)
			continue
		}
		q := req.URL.Query()
		q.Add(dataID, dataId)
		q.Add(GROUP, group)
		req.URL.RawQuery = q.Encode()
		if skipContentCache && cacheData != nil {
			if cacheData.GetLastModifiedHeader() != "" {
				req.Header.Set(IfModifiedSince, cacheData.GetLastModifiedHeader())
			}
			if cacheData.MD5() != "" {
				req.Header.Set(ContentMd5, cacheData.MD5())
			}
		}
		req.Header.Set(AcceptEncoding, "gzip,deflate")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("未知异常", err)
			s.rotateToNextDomain()
			continue
		}
		statusCode := resp.StatusCode
		switch statusCode {
		case http.StatusOK:
			return s.getSuccess(dataId, group, cacheData, resp)
		case http.StatusNotModified:
			return s.getNotModified(dataId, cacheData, resp)
		case http.StatusNotFound:
			log.Println("没有找到DataID为:" + dataId + "对应的配置信息")
			cacheData.SetMD5("")
			s.snapshotConfigInfoProcessor.RemoveSnapshot(dataId, group)
			return "", nil
		case http.StatusServiceUnavailable:
			s.rotateToNextDomain()
			break
		default:
			log.Println("获取修改过的DataID列表的请求回应的HTTP State: ", statusCode)
			s.rotateToNextDomain()
		}
		resp.Body.Close()

	}

	pos := atomic.LoadInt64(&s.domainNamePos)
	domainName, _ := s.diamondConfigure.GetDomainNameList().Get(int(pos))
	return "", fmt.Errorf("获取修改过的DataID列表超时:%v, 超时时间为:%v", domainName, timeout)
}

/**
 * 回馈的结果为RP_NO_CHANGE，则整个流程为：<br>
 * 1.检查缓存中的MD5码与返回的MD5码是否一致，如果不一致，则删除缓存行。重新再次查询。<br>
 * 2.如果MD5码一致，则直接返回NULL<br>
 */
func (s *Subscriber) getNotModified(dataId string, cacheData *configinfo.CacheData, resp *http.Response) (string, error) {
	header := resp.Header
	md5 := header.Get(ContentMd5)
	if md5 == "" {
		return "", errors.New("RP_NO_CHANGE返回的结果中没有MD5码")
	}
	if cacheData.MD5() != md5 {
		lastMd5 := cacheData.MD5()
		cacheData.SetMD5("")
		cacheData.SetLastModifiedHeader("")
		return "", fmt.Errorf("MD5码校验对比出错,DataID为:[%v],上次MD5为:[%v],本次MD5为:[%v]", dataId, lastMd5, md5)
	}

	cacheData.SetMD5(md5)
	s.changeSpacingInterval(header)
	log.Println("DataID: " + dataId + ", 对应的configInfo没有变化")
	return "", nil
}

/**
 * 回馈的结果为RP_OK，则整个流程为：<br>
 * 1.获取配置信息，如果配置信息为空或者抛出异常，则抛出运行时异常<br>
 * 2.检测配置信息是否符合回馈结果中的MD5码，不符合，则再次获取配置信息，并记录日志<br>
 * 3.符合，则存储LastModified信息和MD5码，调整查询的间隔时间，将获取的配置信息发送给客户的监听器<br>
 */
func (s *Subscriber) getSuccess(dataId, group string, cacheData *configinfo.CacheData, resp *http.Response) (string, error) {
	configInfo := getContent(resp)
	if configInfo == "" {
		return "", errors.New("RP_OK获取了错误的配置信息")
	}
	header := resp.Header
	md5 := header.Get(ContentMd5)
	if md5 == "" {
		return "", errors.New("RP_OK返回的结果中没有MD5码, " + configInfo)
	}
	if !checkContent(configInfo, md5) {
		return "", fmt.Errorf("配置信息的MD5码校验出错,DataID为:[%s]配置信息为:[%s]MD5为:[%s]", dataId, configInfo, md5)
	}
	lastModified := header.Get(LastModified)
	if lastModified == "" {
		return "", errors.New("RP_OK返回的结果中没有lastModifiedHeader")
	}

	cacheData.SetLastModifiedHeader(lastModified)
	cacheData.SetMD5(md5)

	s.changeSpacingInterval(header)
	// 设置到本地cache
	key := makeCacheKey(dataId, group)
	s.contentCache.Put(key, configInfo)
	// 记录接收到的数据
	builder := strings.Builder{}
	builder.WriteString("dataId=")
	builder.WriteString(dataId)
	builder.WriteString(" ,group=")
	builder.WriteString(group)
	builder.WriteString(" ,content=")
	builder.WriteString(configInfo)
	log.Println(builder.String())

	return configInfo, nil
}

/**
 * 设置新的消息轮询间隔时间
 */
func (s *Subscriber) changeSpacingInterval(header http.Header) {
	spacingIntervalHeaders := header.Get(SpacingInterval)
	interval, err := strconv.Atoi(spacingIntervalHeaders)
	if err != nil {
		s.diamondConfigure.SetPollingIntervalTime(int64(interval))
	} else {
		log.Println("设置下次间隔时间失败", err)
	}
}

func (s *Subscriber) getCacheData(dataId, group string) *configinfo.CacheData {
	value, ok := s.cache.Load(dataId)
	if ok && value != nil {
		cacheDatas, _ := value.(sync.Map)
		val, ok := cacheDatas.Load(group)
		if ok && val != nil {
			cacheData, _ := val.(*configinfo.CacheData)
			return cacheData
		}
	}

	cacheData := configinfo.NewCacheData(dataId, group)
	var newCache sync.Map
	var oldCache sync.Map
	actual, loaded := s.cache.LoadOrStore(dataId, newCache)
	if !loaded || actual == nil {
		oldCache = actual.(sync.Map)
		oldCache = newCache
	}
	act, loaded := oldCache.LoadOrStore(group, cacheData)
	if loaded && act != nil {
		v, _ := oldCache.Load(group)
		tmp, _ := v.(configinfo.CacheData)
		cacheData = &tmp
	}
	return cacheData
}

func (s *Subscriber) checkLocalConfigInfo() {
	s.cache.Range(func(key, value interface{}) bool {
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
			cacheData, ok := value.(*configinfo.CacheData)
			if !ok {
				return true
			}
			configInfo, err := s.getLocalConfigureInfomation(cacheData)
			if err != nil {
				log.Println("向本地索要配置信息的过程抛异常", err)
				return true
			}
			if configInfo != "" {
				log.Println("本地配置信息被读取, dataId:" + cacheData.DataId() + ", group:" + cacheData.Group())
				s.popConfigInfo(cacheData, configInfo)
				return true
			}
			if cacheData.UseLocalConfigInfo() {
				return true
			}
			return true
		})
		return true
	})
}

/**
 * 将订阅信息抛给客户的监听器
 *
 */
func (s *Subscriber) popConfigInfo(cacheData *configinfo.CacheData, configInfo string) {
	configureInfomation := configinfo.NewConfigureInformation()
	configureInfomation.DataId = cacheData.DataId()
	configureInfomation.Group = cacheData.Group()
	configureInfomation.ConfigureInfo = configInfo
	cacheData.IncrementFetchCountAndGet()
	s.subscriberListener.ReceiveConfigInfo(configureInfomation)
	s.saveSnapshot(cacheData.DataId(), cacheData.Group(), configInfo)
}

func (s *Subscriber) getLocalConfigureInfomation(cacheData *configinfo.CacheData) (string, error) {
	if !s.isRun {
		return "", errors.New("DiamondSubscriber不在运行状态中，无法获取本地ConfigureInfomation")
	}
	return s.localConfigInfoProcessor.GetLocalConfigureInformation(cacheData, false)
}

func (s *Subscriber) checkUpdateDataIds(timeout int) (*hashset.Set, error) {
	if !s.isRun {
		return nil, errors.New("subscriber不在运行状态中，无法获取修改过的DataID列表")
	}

	probeUpdateString := s.getProbeUpdateString()
	if probeUpdateString == "" {
		return nil, errors.New("getProbeUpdateString is empty")
	}
	waitTime := 0
	for 0 == timeout || timeout > waitTime {
		onceTimeOut := s.getOnceTimeOut(waitTime, timeout)
		waitTime += onceTimeOut

		client := &http.Client{Timeout: time.Duration(onceTimeOut) * time.Millisecond}
		pos := atomic.LoadInt64(&s.domainNamePos)
		domainName, _ := s.diamondConfigure.GetDomainNameList().Get(int(pos))
		domainNamePort := urlutil.GetURL(domainName.(string), s.diamondConfigure.GetPort(), GetProbeModifyUrl)

		params := url.Values{}
		params.Set(probeModifyRequest, probeUpdateString)
		req, err := http.NewRequest("POST", domainNamePort, strings.NewReader(params.Encode()))
		if err != nil {
			log.Println("Can't NewRequest", err)
			continue
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			log.Println("未知异常", err)
			continue
		}
		statusCode := resp.StatusCode
		switch statusCode {
		case http.StatusOK:
			return getUpdateDataIds(resp), nil
		case http.StatusServiceUnavailable:
			s.rotateToNextDomain()
			break
		default:
			log.Println("获取修改过的DataID列表的请求回应的HTTP State: ", statusCode)
			s.rotateToNextDomain()
		}
		resp.Body.Close()
	}

	pos := atomic.LoadInt64(&s.domainNamePos)
	domainName, _ := s.diamondConfigure.GetDomainNameList().Get(int(pos))
	return nil, fmt.Errorf("获取修改过的DataID列表超时:%v, 超时时间为:%v", domainName, timeout)
}

func (s *Subscriber) getProbeUpdateString() string {
	probeModifyBuilder := strings.Builder{}
	s.cache.Range(func(key, value interface{}) bool {
		dataId, _ := key.(string)
		if value == nil {
			return true
		}
		groupCache, ok := value.(sync.Map)
		if !ok {
			return true
		}

		l := maputil.LengthOfSyncMap(groupCache)
		groupCache.Range(func(key, value interface{}) bool {
			if value == nil {
				return true
			}
			data, ok := value.(*configinfo.CacheData)
			// 非使用本地配置，才去diamond server检查
			if !ok {
				return true
			}
			if !data.UseLocalConfigInfo() {
				probeModifyBuilder.WriteString(dataId)
				probeModifyBuilder.WriteString(wordSeparator)

				if data.Group() != "" {
					probeModifyBuilder.WriteString(data.Group())
					probeModifyBuilder.WriteString(wordSeparator)
				} else {
					probeModifyBuilder.WriteString(wordSeparator)
				}

				if data.MD5() != "" {
					probeModifyBuilder.WriteString(data.MD5())
				}

				if l > 1 {
					probeModifyBuilder.WriteString(lineSeparator)
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

func (s *Subscriber) rotateToNextDomain() {
	domainNameList := s.diamondConfigure.GetDomainNameList()
	domainNameCount := domainNameList.Size()
	index := atomic.AddInt64(&s.domainNamePos, 1)
	if index < 0 {
		index = -index
	}
	if domainNameCount == 0 {
		log.Println("diamond server address is empty, please concat admin")
		return
	}
	atomic.StoreInt64(&s.domainNamePos, index%int64(domainNameCount))
	if domainNameList.Size() > 0 {
		domainName, _ := domainNameList.Get(int(atomic.LoadInt64(&s.domainNamePos)))
		log.Println("rotateTo doaminName：", domainName)
	}
}

//random domain name pos
func (s *Subscriber) randomDomainNamePos() {
	rand.Seed(time.Now().UnixNano())
	domainList := s.diamondConfigure.GetDomainNameList()
	if !domainList.Empty() {
		s.domainNamePos = rand.Int63n(int64(domainList.Size()))
	}
}

//saveSnapshot save config data to snapshot local file
func (s *Subscriber) saveSnapshot(dataId, group, config string) {
	if config != "" {
		err := s.snapshotConfigInfoProcessor.SaveSnaptshot(dataId, group, config)
		if err != nil {
			log.Println("保存snapshot出错,dataId="+dataId+",group="+group, err)
		}
	}
}

func (s *Subscriber) checkSnapshot() {
	s.cache.Range(func(key, value interface{}) bool {
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
			cacheData, ok := value.(*configinfo.CacheData)
			if !ok {
				return true
			}
			//if get local file config failed and get config from diamond server all failed
			//then get config from last time succeed local snapshot
			if !cacheData.UseLocalConfigInfo() && cacheData.GetFetchCount() == 0 {
				configInfo := s.getSnapshotConfiginfomation(cacheData.DataId(), cacheData.Group())
				if configInfo != "" {
					s.popConfigInfo(cacheData, configInfo)
				}
			}
			return true
		})
		return true
	})
}

func (s *Subscriber) getSnapshotConfiginfomation(dataId, group string) string {
	if group == "" {
		group = configinfo.DefaultGroup
	}
	cacheData := s.getCacheData(dataId, group)
	config, err := s.snapshotConfigInfoProcessor.GetConfigInfomation(dataId, group)
	if err != nil {
		log.Println("get snapshot failed， dataId="+dataId+",group="+group, err)
		return ""
	}
	if config != "" && cacheData != nil {
		cacheData.IncrementFetchCountAndGet()
	}
	return config
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
		log.Println("reader.Read modifiedDataIdsString faield", err)
	}

	if modifiedDataIdsString != "" {
		if strings.HasPrefix(modifiedDataIdsString, "OK") {
			log.Println("result:" + modifiedDataIdsString)
		} else {
			log.Println("data change:" + modifiedDataIdsString)
		}
	}

	modifiedDataIdSet := hashset.New()
	modifiedDataIdStrings := strings.Split(modifiedDataIdsString, lineSeparator)
	for _, modifiedDataIdString := range modifiedDataIdStrings {
		if modifiedDataIdString != "" && modifiedDataIdsString != "OK" {
			modifiedDataIdSet.Add(modifiedDataIdString)
		}
	}
	return modifiedDataIdSet
}

func checkContent(configInfo, md5 string) bool {
	realMd5 := stringutil.GetMd5(configInfo)
	if realMd5 == "" {
		return md5 == ""
	}
	return realMd5 == md5
}

func getContent(resp *http.Response) string {
	contentBuilder := strings.Builder{}

	if isZipContent(resp.Header) {
		//unzip
		var buf [1024 * 4]byte
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Println("unzip failed", err)
		} else {
			for {
				n, err := reader.Read(buf[:])
				if err != nil || n == 0 {
					log.Println("unzip failed or finished", err)
					break
				}
				contentBuilder.WriteString(string(buf[:n]))
			}
			reader.Close()
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("read resp.Body failed", err)
		}
		content := string(body)
		if content == "" {
			return ""
		}
		contentBuilder.WriteString(content)
	}
	return contentBuilder.String()
}

func isZipContent(header http.Header) bool {
	acceptEncoding := header.Get(ContentEncoding)
	if acceptEncoding != "" {
		if strings.Index(strings.ToLower(acceptEncoding), "gzip") > -1 {
			return true
		}
	}
	return false
}

//concat cache key
func makeCacheKey(dataId, group string) string {
	key := dataId + "-" + group
	return key
}
