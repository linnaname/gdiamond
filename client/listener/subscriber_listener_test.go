package listener

import (
	"fmt"
	"gdiamond/client/internal/configinfo"
	"gdiamond/util/maputil"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type _S struct {
	suite.Suite
	dsl *DefaultSubscriberListener
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	dl := NewDefaultSubscriberListener()
	s.dsl = &dl
}

func (s *_S) TearDownSuite() {
	fmt.Printf("TearDownSuite() ...\n")
}

func TestSS(t *testing.T) {
	suite.Run(t, new(_S))
}

type TestBMangerListener struct {
}

func (tml TestBMangerListener) ReceiveConfigInfo(configInfo string) {
	fmt.Println("B", configInfo)
}

type TestAMangerListener struct {
}

func (tml TestAMangerListener) ReceiveConfigInfo(configInfo string) {
	fmt.Println("A", configInfo)
}

func (s *_S) TestDefaultSubscriberListener_AddManagerListeners() {
	listeners := singlylinkedlist.New()
	listeners.Add(TestBMangerListener{})
	s.dsl.AddManagerListeners("linnana", "DEFAULT_GROUP", listeners)
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.dsl.allListeners), int64(1))
}

func (s *_S) TestDefaultSubscriberListener_RemoveManagerListeners() {
	s.dsl.RemoveManagerListeners("linnana", "DEFAULT_GROUP")
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.dsl.allListeners), int64(0))
}

func (s *_S) TestDefaultSubscriberListener_ReceiveConfigInfo() {
	listeners := singlylinkedlist.New()
	listeners.Add(TestAMangerListener{})
	listeners.Add(TestBMangerListener{})
	s.dsl.AddManagerListeners("linnana", "DEFAULT_GROUP", listeners)
	configureInfomation := configinfo.NewConfigureInformation()
	configureInfomation.Group = "DEFAULT_GROUP"
	configureInfomation.DataId = "linnana"
	configureInfomation.ConfigureInfo = "content"
	s.dsl.ReceiveConfigInfo(configureInfomation)
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.dsl.allListeners), int64(1))

	key := makeKey("linnana", "DEFAULT_GROUP")
	value, ok := s.dsl.allListeners.Load(key)
	assert.True(s.T(), ok)
	assert.NotNil(s.T(), value)
}

func TestName(t *testing.T) {
	var sm sync.Map
	assert.NotNil(t, sm)
	sl := singlylinkedlist.New()
	sl.Add("once")
	sm.Store("list", sl)
	sl.Add("more")
	fmt.Println(sl.Size())
	value, ok := sm.Load("list")
	assert.True(t, ok)
	assert.NotNil(t, value)
	l, ok := value.(*singlylinkedlist.List)
	fmt.Println(ok)
	fmt.Println(l.Size())
	fmt.Println(l)

}
