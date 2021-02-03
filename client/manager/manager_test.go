package manager

import (
	"fmt"
	"gdiamond/client/listener"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	dataId = "linna"
	group  = "DEFAULT_GROUP"
)

func TestName(t *testing.T) {
	ins := GetSubscriberInstance()
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

func TestNewManager(t *testing.T) {
	dm := NewManager()
	assert.NotNil(t, dm)
}

func TestDefaultManager_GetAndSetManagerListener(t *testing.T) {
	dm := NewManager()
	dm.SetManagerListener(dataId, group, A{})
	assert.NotNil(t, dm.GetManagerListeners())
	assert.Equal(t, dm.GetManagerListeners().Size(), 1)
}

func TestDefaultManager_Close(t *testing.T) {
	dm := NewManager()
	dm.Close()
}

func TestDefaultManager_GetConfig(t *testing.T) {
	dm := NewManager()
	content := dm.GetConfig(dataId, group, 1000)
	assert.NotEmpty(t, content)
	fmt.Println(content)
}

func TestDefaultManager_PublishConfig(t *testing.T) {
	dm := NewManager()
	b := dm.PublishConfig("linna3", group, "test listener999")
	assert.True(t, b)
}

func TestDefaultManager_GetConfigAndSetListener(t *testing.T) {
	dm := NewManager()
	content := dm.GetConfigAndSetListener("linna3", group, 1000, A{})
	assert.NotEmpty(t, content)
	time.Sleep(time.Minute * 10)
}
