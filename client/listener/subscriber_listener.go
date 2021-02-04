package listener

import (
	"errors"
	"gdiamond/client/internal/configinfo"
	"gdiamond/client/internal/logger"
	"gdiamond/util/maputil"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/sirupsen/logrus"
	"sync"
)

//SubscriberListener listener
type SubscriberListener interface {
	//implement this method to handle notify mananger listener ManagerListener
	ReceiveConfigInfo(configureInfomation *configinfo.ConfigureInformation)
}

//DefaultSubscriberListener default listener
type DefaultSubscriberListener struct {
	// k:dataId + group  v:listeners  v is list
	allListeners sync.Map
}

//NewDefaultSubscriberListener new
func NewDefaultSubscriberListener() DefaultSubscriberListener {
	dl := DefaultSubscriberListener{}
	return dl
}

//ReceiveConfigInfo notify all listener which implement ManagerListener
//it's a pointer receiver,   pointer receiver can't invoke   ReceiveConfigInfo in subsriber
func (d DefaultSubscriberListener) ReceiveConfigInfo(configureInfomation *configinfo.ConfigureInformation) {
	dataId := configureInfomation.DataId
	group := configureInfomation.Group

	if dataId == "" {
		logger.Logger.WithFields(logrus.Fields{}).Error("receiveConfigInfo dataId can't  be empty")
		return
	}
	key := makeKey(dataId, group)
	value, ok := d.allListeners.Load(key)
	if !ok || value == nil {
		if dataId == "" {
			logger.Logger.WithFields(logrus.Fields{
				"dataId": dataId,
				"group":  group,
			}).Info("[notify-listener] no listener")
			return
		}
		return
	}
	listeners, _ := value.(*singlylinkedlist.List)
	listeners.Each(func(index int, value interface{}) {
		listener, _ := value.(ManagerListener)
		err := notifyListener(configureInfomation, listener)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"dataId": dataId,
				"group":  group,
			}).Info("[notify-listener] call listener failed")
		}
	})
}

//AddManagerListeners If dataId or addListeners is empty it do nothing,if group is empty it will be assign to DEFAULT_GROUP
func (d *DefaultSubscriberListener) AddManagerListeners(dataId, group string, addListeners *singlylinkedlist.List) {
	if dataId == "" || addListeners.Empty() {
		return
	}

	key := makeKey(dataId, group)
	value, ok := d.allListeners.Load(key)

	listenerList := singlylinkedlist.New()
	if !ok || value == nil {
		actual, loaded := d.allListeners.LoadOrStore(key, listenerList)
		if actual != nil && loaded {
			oldList, _ := actual.(*singlylinkedlist.List)
			listenerList = oldList
		}
	}
	addListeners.Each(func(index int, value interface{}) {
		listenerList.Add(value)
	})
}

//RemoveManagerListeners remove all listener of dataId/group  If dataId is empty it do nothing,if group is empty it will be assign to DEFAULT_GROUP
func (d *DefaultSubscriberListener) RemoveManagerListeners(dataId, group string) {
	if dataId == "" {
		return
	}
	key := makeKey(dataId, group)
	d.allListeners.Delete(key)
}

func (d *DefaultSubscriberListener) ClearManagerListeners() {
	maputil.ClearSyncMap(d.allListeners)
}

//notify listener async
func notifyListener(configureInfomation *configinfo.ConfigureInformation, listener ManagerListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}

	dataId := configureInfomation.DataId
	group := configureInfomation.Group
	content := configureInfomation.ConfigureInfo
	logger.Logger.WithFields(logrus.Fields{
		"dataId":  dataId,
		"group":   group,
		"content": content,
	}).Debug("[notify-listener] calling listener")
	//notify listener async
	go listener.ReceiveConfigInfo(content)

	return nil
}

//concat dataId and group to key
func makeKey(dataId, group string) string {
	if group == "" {
		group = configinfo.DefaultGroup
	}
	return dataId + "_" + group
}
