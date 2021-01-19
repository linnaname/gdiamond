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
	println(configInfo)
}

func TestNewManager(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetConfigureInfomation(1000)
	//assert.NotEmpty(t, content)
	println(content)
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
	content := dm.GetConfigureInfomation(1000)
	fmt.Println(content)
}

func TestDefaultManager_GetAvailableConfigureInfomation(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetAvailableConfigureInfomation(1000)
	fmt.Println(content)
}

func TestDefaultManager_GetAvailableConfigureInfomationFromSnapshot(t *testing.T) {
	dm := NewManager("linna", "DEFAULT_GROUP", A{})
	content := dm.GetAvailableConfigureInfomationFromSnapshot(1000)
	fmt.Println(content)
}
