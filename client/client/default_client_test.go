package client

import (
	"fmt"
	"gdiamond/client/listener"
	"gdiamond/client/subscriber"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	dataId = "linna"
	group  = "DEFAULT_GROUP"
)

func TestName(t *testing.T) {
	ins := subscriber.GetSubscriberInstance()
	sl := ins.GetSubscriberListener()
	dl, ok := sl.(listener.DefaultSubscriberListener)
	assert.True(t, ok)
	assert.NotNil(t, dl)
}

type A struct {
}

func (a A) ReceiveConfigInfo(configInfo string) {
	println("ReceiveConfigInfo:", configInfo)
}

func TestNewClient(t *testing.T) {
	dm := NewClient()
	assert.NotNil(t, dm)
}

func TestDefaultClient_GetAndSetManagerListener(t *testing.T) {
	dm := NewClient()
	dm.SetManagerListener(dataId, group, A{})
	assert.NotNil(t, dm.GetManagerListeners())
	assert.Equal(t, dm.GetManagerListeners().Size(), 1)
}

func TestDefaultClient_Close(t *testing.T) {
	dm := NewClient()
	dm.Close()
}

func TestDefaultClient_GetConfig(t *testing.T) {
	dm := NewClient()
	content := dm.GetConfig(dataId, group, 1000)
	assert.NotEmpty(t, content)
	fmt.Println(content)
}

func TestDefaultClient_PublishConfig(t *testing.T) {
	dm := NewClient()
	b := dm.PublishConfig("linna3", group, "test publish22")
	assert.True(t, b)
}

func TestDefaultClient_GetConfigAndSetListener(t *testing.T) {
	dm := NewClient()
	content := dm.GetConfigAndSetListener("linna3", group, 1000, A{})
	assert.NotEmpty(t, content)
	time.Sleep(time.Minute * 20)
}
