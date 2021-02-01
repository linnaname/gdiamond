package manager

import (
	"fmt"
	"gdiamond/client/listener"
	"github.com/stretchr/testify/assert"
	"testing"
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
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	assert.NotNil(t, dm)
}

func TestDefaultManager_GetAndSetManagerListener(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	dm.SetManagerListener(A{})
	assert.NotNil(t, dm.GetManagerListeners())
	assert.Equal(t, dm.GetManagerListeners().Size(), 1)
}

func TestDefaultManager_Close(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	dm.Close()
	assert.NotNil(t, dm.GetManagerListeners())
	assert.Equal(t, dm.GetManagerListeners().Size(), 1)
}

func TestDefaultManager_GetConfigureInfomation(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetConfigureInformation(1000)
	assert.NotEmpty(t, content)
	fmt.Println(content)
}

func TestDefaultManager_GetAvailableConfigureInfomation(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetAvailableConfigureInformation(1000)
	assert.NotEmpty(t, content)
	fmt.Println(content)
}

func TestDefaultManager_GetAvailableConfigureInfomationFromSnapshot(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetAvailableConfigureInformationFromSnapshot(1000)
	assert.NotEmpty(t, content)
	fmt.Println(content)
}

func TestDefaultManager_PublishConfig(t *testing.T) {
	dm := NewManager("my.test", "DEFAULT_GROUP", A{})
	b := dm.PublishConfig("whaterver it's")
	assert.True(t, b)
}

func TestDefaultManager_ManagerListener(t *testing.T) {
	NewManager("linna", "DEFAULT_GROUP", A{})
	select {}
}
