package manager

import (
	"gdiamond/client/listener"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

//Manager user api
type Manager interface {
	GetConfig(dataId, group string, timeout int) string

	GetConfigAndSetListener(dataId, group string, timeout int, mListener listener.ManagerListener) string

	PublishConfig(dataId, group string, content string) bool

	////GetConfigureInformation 同步获取配置信息,,此方法优先从gdiamond/data/config-data/${group}/${dataId}
	////下获取配置文件，如果没有，则从gdiamond server获取配置信息
	//GetConfigureInformation(dataId,group string, timeout int) string
	//
	////GetAvailableConfigureInformation 同步获取一份有效的配置信息，按照local file ->diamond server -> 上一次正确配置的snapshot
	////的优先顺序获取， 如果这些途径都无效，则返回""
	//GetAvailableConfigureInformation(dataId, group string, timeout int) string
	//
	////GetAvailableConfigureInformationFromSnapshot 同步获取一份有效的配置信息，按照上一次正确配置的snapshot->local file-> diamond server
	////的优先顺序获取， 如果这些途径都无效，则返回""
	//GetAvailableConfigureInformationFromSnapshot(dataId, group string, timeout int) string
	//

	//GetManagerListeners 返回该DiamondManager设置的listener列表
	GetManagerListeners() *singlylinkedlist.List

	//SetManagerListener 设置ManagerListener，每当收到一个DataID对应的配置信息，则客户设置的ManagerListener会接收到这个配置信息
	SetManagerListener(dataId, group string, mListener listener.ManagerListener)

	//Close close
	Close()
}

//DefaultManager default manager
type DefaultManager struct {
	subscriber       *Subscriber
	managerListeners *singlylinkedlist.List
	//dataId           string
	//group            string
}

//NewManager new
func NewManager() *DefaultManager {
	dm := &DefaultManager{}
	dm.subscriber = GetSubscriberInstance()
	if dm.managerListeners == nil {
		dm.managerListeners = singlylinkedlist.New()
	}
	dm.subscriber.Start()
	return dm
}

//Close implement method
func (dm *DefaultManager) Close() {
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.ClearManagerListeners()
	dm.subscriber.ClearCache()
	if dm.subscriber.GetDataIds().Size() == 0 {
		dm.subscriber.Close()
	}
}

func (dm *DefaultManager) GetConfig(dataId, group string, timeout int) string {
	dm.subscriber.AddDataId(dataId, group)
	return dm.subscriber.GetConfigureInformation(dataId, group, timeout)
}

func (dm *DefaultManager) GetConfigAndSetListener(dataId, group string, timeout int, mListener listener.ManagerListener) string {
	dm.managerListeners.Clear()
	dm.managerListeners.Add(mListener)
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dataId, group)
	dsl.AddManagerListeners(dataId, group, dm.managerListeners)
	dm.subscriber.SetSubscriberListener(dsl)
	return dm.GetConfig(dataId, group, timeout)
}

//GetManagerListeners  implement method
func (dm *DefaultManager) GetManagerListeners() *singlylinkedlist.List {
	return dm.managerListeners
}

//SetManagerListener  implement method
func (dm *DefaultManager) SetManagerListener(dataId, group string, mListener listener.ManagerListener) {
	dm.managerListeners.Clear()
	dm.managerListeners.Add(mListener)
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dataId, group)
	dsl.AddManagerListeners(dataId, group, dm.managerListeners)
	dm.subscriber.SetSubscriberListener(dsl)
}

//PublishConfig publish config info
func (dm *DefaultManager) PublishConfig(dataId, group string, content string) bool {
	err := dm.subscriber.PublishConfigureInformation(dataId, group, content)
	if err != nil {
		return false
	}
	return true
}

//GetConfigureInformation implement method
func (dm *DefaultManager) getConfigureInformation(dataId, group string, timeout int) string {
	return dm.subscriber.GetConfigureInformation(dataId, group, timeout)
}

//GetAvailableConfigureInformation implement method
func (dm *DefaultManager) getAvailableConfigureInformation(dataId, group string, timeout int) string {
	return dm.subscriber.GetAvailableConfigureInformation(dataId, group, timeout)
}

//GetAvailableConfigureInformationFromSnapshot  implement method
func (dm *DefaultManager) getAvailableConfigureInformationFromSnapshot(dataId, group string, timeout int) string {
	return dm.subscriber.GetAvailableConfigureInformationFromSnapshot(dataId, group, timeout)
}
