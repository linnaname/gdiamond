package manager

import (
	"gdiamond/client/listener"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

//Manager user api
type Manager interface {

	//GetConfigureInformation 同步获取配置信息,,此方法优先从gdiamond/data/config-data/${group}/${dataId}
	//下获取配置文件，如果没有，则从gdiamond server获取配置信息
	GetConfigureInformation(timeout int) string

	//GetAvailableConfigureInformation 同步获取一份有效的配置信息，按照local file ->diamond server -> 上一次正确配置的snapshot
	//的优先顺序获取， 如果这些途径都无效，则返回""
	GetAvailableConfigureInformation(timeout int) string

	//GetAvailableConfigureInformationFromSnapshot 同步获取一份有效的配置信息，按照上一次正确配置的snapshot->local file-> diamond server
	//的优先顺序获取， 如果这些途径都无效，则返回""
	GetAvailableConfigureInformationFromSnapshot(timeout int) string

	//GetManagerListeners 返回该DiamondManager设置的listener列表
	GetManagerListeners() *singlylinkedlist.List

	//SetManagerListener 设置ManagerListener，每当收到一个DataID对应的配置信息，则客户设置的ManagerListener会接收到这个配置信息
	SetManagerListener(mListener listener.ManagerListener)

	//Close close
	Close()
}

//DefaultManager default manager
type DefaultManager struct {
	subscriber       *Subscriber
	managerListeners *singlylinkedlist.List
	dataId           string
	group            string
}

//NewManager new
func NewManager(dataId, group string, mlistener listener.ManagerListener) *DefaultManager {
	dm := &DefaultManager{dataId: dataId, group: group}
	dm.subscriber = GetSubscriberInstance()
	if dm.managerListeners == nil {
		dm.managerListeners = singlylinkedlist.New()
	}
	dm.managerListeners.Add(mlistener)

	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.AddManagerListeners(dataId, group, dm.managerListeners)
	dm.subscriber.SetSubscriberListener(dsl)
	dm.subscriber.AddDataId(dataId, group)
	dm.subscriber.Start()
	return dm
}

//Close implement method
func (dm *DefaultManager) Close() {
	/**
	 * 因为同一个DataID只能对应一个MnanagerListener，所以，关闭时一次性关闭所有ManagerListener即可
	 */
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dm.dataId, dm.group)

	dm.subscriber.RemoveDataId(dm.dataId, dm.group)
	if dm.subscriber.GetDataIds().Size() == 0 {
		dm.subscriber.Close()
	}

}

//GetConfigureInformation implement method
func (dm *DefaultManager) GetConfigureInformation(timeout int) string {
	return dm.subscriber.GetConfigureInformation(dm.dataId, dm.group, timeout)
}

//GetAvailableConfigureInformation implement method
func (dm *DefaultManager) GetAvailableConfigureInformation(timeout int) string {
	return dm.subscriber.GetAvailableConfigureInformation(dm.dataId, dm.group, timeout)
}

//GetAvailableConfigureInformationFromSnapshot  implement method
func (dm *DefaultManager) GetAvailableConfigureInformationFromSnapshot(timeout int) string {
	return dm.subscriber.GetAvailableConfigureInformationFromSnapshot(dm.dataId, dm.group, timeout)
}

//GetManagerListeners  implement method
func (dm *DefaultManager) GetManagerListeners() *singlylinkedlist.List {
	return dm.managerListeners
}

//SetManagerListener  implement method
func (dm *DefaultManager) SetManagerListener(mListener listener.ManagerListener) {
	dm.managerListeners.Clear()
	dm.managerListeners.Add(mListener)
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dm.dataId, dm.group)
	dsl.AddManagerListeners(dm.dataId, dm.group, dm.managerListeners)
}

//PublishConfig publish config info
func (dm *DefaultManager) PublishConfig(content string) bool {
	err := dm.subscriber.PublishConfigureInformation(dm.dataId, dm.group, content)
	if err != nil {
		return false
	}
	return true
}
