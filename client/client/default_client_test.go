package client

import (
	"fmt"
	"gdiamond/client/listener"
	"gdiamond/client/subscriber"
	"github.com/stretchr/testify/assert"
	"strconv"
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
	b := dm.PublishConfig("linna3", group, "test publish442221")
	assert.True(t, b)
}

func TestDefaultClient_GetConfigAndSetListener(t *testing.T) {
	dm := NewClient()
	content := dm.GetConfigAndSetListener("linna3", group, 1000, A{})
	assert.NotEmpty(t, content)
	time.Sleep(time.Minute * 20)
}

func BenchmarkGetConfig(b *testing.B) {
	dm := NewClient()
	for i := 0; i < b.N; i++ {
		dm.GetConfig("linna3"+strconv.Itoa(i), group, 1000)
	}
}

func BenchmarkPublishConfig(b *testing.B) {
	dm := NewClient()
	for i := 0; i < b.N; i++ {
		dm.PublishConfig("linna3"+strconv.Itoa(i), group, "test publish442221"+strconv.Itoa(i))
	}
}

func BenchmarkParallelGetConfig(b *testing.B) {
	// 测试一个对象或者函数在多线程的场景下面是否安全
	dm := NewClient()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			dm.GetConfig(dataId, group, 1000)
		}
	})
}
