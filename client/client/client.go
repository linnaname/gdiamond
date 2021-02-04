package client

import (
	"gdiamond/client/listener"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
)

//Client user api
type Client interface {

	//GetConfig get config  by dataId and group
	//if timeout or internal error it return empty string
	GetConfig(dataId, group string, timeout int) string

	//GetConfigAndSetListener get config  by dataId and group
	//if timeout or internal error it return empty string
	//need to implement listener.ManagerListener to receive config changed
	GetConfigAndSetListener(dataId, group string, timeout int, mListener listener.ManagerListener) string

	//PublishConfig publish or update  config by dataId and group
	//if internal error it will return false
	PublishConfig(dataId, group string, content string) bool

	//GetManagerListeners get all listener
	GetManagerListeners() *singlylinkedlist.List

	//SetManagerListener set listener by dataId and group
	SetManagerListener(dataId, group string, mListener listener.ManagerListener)

	//Close release resource
	Close()
}
