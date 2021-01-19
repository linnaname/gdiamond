package listener

import (
	"errors"
	"gdiamond/client/configinfo"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"log"
	"sync"
)

type SubscriberListener interface {
	//implement this method to handle notify mananger listener ManagerListener
	ReceiveConfigInfo(configureInfomation *configinfo.ConfigureInfomation)
}

type DefaultSubscriberListener struct {
	// k:dataId + group  v:listeners  v is list
	allListeners sync.Map
}

func NewDefaultSubscriberListener() DefaultSubscriberListener {
	dl := DefaultSubscriberListener{}
	return dl
}

//ReceiveConfigInfo notify all listener which implement ManagerListener
//it's a pointer receiver,   pointer receiver can't invoke   ReceiveConfigInfo in subsriber
func (d DefaultSubscriberListener) ReceiveConfigInfo(configureInfomation *configinfo.ConfigureInfomation) {
	dataId := configureInfomation.DataId
	group := configureInfomation.Group

	if dataId == "" {
		log.Println("[receiveConfigInfo] dataId is null")
		return
	}
	key := makeKey(dataId, group)
	value, ok := d.allListeners.Load(key)
	if !ok || value == nil {
		log.Println("[notify-listener] no listener for dataId=" + dataId + ", group=" + group)
		return
	}
	listeners, _ := value.(*singlylinkedlist.List)
	listeners.Each(func(index int, value interface{}) {
		listener, _ := value.(ManagerListener)
		err := notifyListener(configureInfomation, listener)
		if err != nil {
			log.Println("call listener error, dataId="+dataId+", group="+group, err)
		}
	})
}

/**
*AddManagerListeners
 If dataId or addListeners is empty it do nothing,if group is empty it will be assign to DEFAULT_GROUP
*/
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

/**
* RemoveManagerListeners remove all listener of dataId/group
  If dataId is empty it do nothing,if group is empty it will be assign to DEFAULT_GROUP
*/
func (d *DefaultSubscriberListener) RemoveManagerListeners(dataId, group string) {
	if dataId == "" {
		return
	}
	key := makeKey(dataId, group)
	d.allListeners.Delete(key)
}

//notify listener async
func notifyListener(configureInfomation *configinfo.ConfigureInfomation, listener ManagerListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}

	dataId := configureInfomation.DataId
	group := configureInfomation.Group
	content := configureInfomation.ConfigureInfo
	log.Println("[notify-listener] call listener  for " + dataId + ", " + group + ", " + content)
	//notify listener async
	go listener.ReceiveConfigInfo(content)

	return nil
}

//concat dataId and group to key
func makeKey(dataId, group string) string {
	if group == "" {
		group = configinfo.DEFAULT_GROUP
	}
	return dataId + "_" + group
}
