package manager

import (
	"gdiamond/client/listener"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

type Manager interface {
	/**
	 * 同步获取配置信息,,此方法优先从${user.home}/diamond/data/config-data/${group}/${dataId}
	下获取配置文件，如果没有，则从diamond server获取配置信息
	*/
	GetConfigureInfomation(timeout int) string

	/**
	 * 同步获取一份有效的配置信息，按照<strong>本地文件->diamond服务器->上一次正确配置的snapshot</strong>
	 * 的优先顺序获取， 如果这些途径都无效，则返回""
	 *
	 */
	GetAvailableConfigureInfomation(timeout int) string

	/**
	 * 同步获取一份有效的配置信息，按照<strong>上一次正确配置的snapshot->本地文件->diamond服务器</strong>
	 * 的优先顺序获取， 如果这些途径都无效，则返回""
	 */
	GetAvailableConfigureInfomationFromSnapshot(timeout int) string

	/**
	 * 返回该DiamondManager设置的listener列表
	 */
	GetManagerListeners() *singlylinkedlist.List

	/**
	 * 设置ManagerListener，每当收到一个DataID对应的配置信息，则客户设置的ManagerListener会接收到这个配置信息
	 */
	SetManagerListener(mlistener listener.ManagerListener)

	/**
	关闭
	*/
	Close()
}

type DefaultManager struct {
	subscriber       *Subscriber
	managerListeners *singlylinkedlist.List
	dataId           string
	group            string
}

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
	dm.subscriber.AddDataId(dataId, group)
	dm.subscriber.Start()
	return dm
}

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

func (dm *DefaultManager) GetConfigureInfomation(timeout int) string {
	return dm.subscriber.GetConfigureInfomation(dm.dataId, dm.group, timeout)
}

func (dm *DefaultManager) GetAvailableConfigureInfomation(timeout int) string {
	return dm.subscriber.GetAvailableConfigureInfomation(dm.dataId, dm.group, timeout)
}

func (dm *DefaultManager) GetAvailableConfigureInfomationFromSnapshot(timeout int) string {
	return dm.subscriber.GetAvailableConfigureInfomationFromSnapshot(dm.dataId, dm.group, timeout)
}

func (dm *DefaultManager) GetManagerListeners() *singlylinkedlist.List {
	return dm.managerListeners
}

func (dm *DefaultManager) SetManagerListener(mlistener listener.ManagerListener) {
	dm.managerListeners.Clear()
	dm.managerListeners.Add(mlistener)
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dm.dataId, dm.group)
	dsl.AddManagerListeners(dm.dataId, dm.group, dm.managerListeners)
}
