package client

import (
	"gdiamond/client/listener"
	"gdiamond/client/subscriber"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

//DefaultClient default client implement Client
type DefaultClient struct {
	subscriber       *subscriber.Subscriber
	managerListeners *singlylinkedlist.List
}

//NewClient new
func NewClient() *DefaultClient {
	dm := &DefaultClient{}
	dm.subscriber = subscriber.GetSubscriberInstance()
	if dm.managerListeners == nil {
		dm.managerListeners = singlylinkedlist.New()
	}
	dm.subscriber.Start()
	return dm
}

//Close implement method
func (dm *DefaultClient) Close() {
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.ClearManagerListeners()
	dm.subscriber.ClearCache()
	if dm.subscriber.GetDataIds().Size() == 0 {
		dm.subscriber.Close()
	}
}

//GetConfig get config  by dataId and group
//if timeout or internal error it return empty string
//config sequence:local file -> server -> snapshot file
func (dm *DefaultClient) GetConfig(dataId, group string, timeout int) string {
	dm.subscriber.AddDataId(dataId, group)
	return dm.subscriber.GetConfigureInformation(dataId, group, timeout)
}

//GetConfigAndSetListener get config  by dataId and group
//if timeout or internal error it return empty string
//need to implement listener.ManagerListener to receive config changed
//config sequence:local file -> server -> snapshot file
func (dm *DefaultClient) GetConfigAndSetListener(dataId, group string, timeout int, mListener listener.ManagerListener) string {
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
func (dm *DefaultClient) GetManagerListeners() *singlylinkedlist.List {
	return dm.managerListeners
}

//SetManagerListener  implement method
func (dm *DefaultClient) SetManagerListener(dataId, group string, mListener listener.ManagerListener) {
	dm.managerListeners.Clear()
	dm.managerListeners.Add(mListener)
	sl := dm.subscriber.GetSubscriberListener()
	dsl, _ := sl.(listener.DefaultSubscriberListener)
	dsl.RemoveManagerListeners(dataId, group)
	dsl.AddManagerListeners(dataId, group, dm.managerListeners)
	dm.subscriber.SetSubscriberListener(dsl)
}

//PublishConfig publish config info to server
func (dm *DefaultClient) PublishConfig(dataId, group string, content string) bool {
	err := dm.subscriber.PublishConfigureInformation(dataId, group, content)
	if err != nil {
		return false
	}
	return true
}
