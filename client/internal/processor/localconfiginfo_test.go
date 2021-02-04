package processor

import (
	"fmt"
	"gdiamond/client/internal/configinfo"
	"gdiamond/util/fileutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _S struct {
	suite.Suite
	p *LocalConfigInfoProcessor
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	p := NewLocalConfigInfoProcessor()
	s.p = p
	p.Start("/Users/goranka/gdiamond/data")
}

func (s *_S) TearDownSuite() {
	s.p.Stop()
	fmt.Printf("TearDownSuite() ...\n")
}

func (s *_S) TestStart() {
	assert.True(s.T(), s.p.isRun)
	assert.True(s.T(), fileutil.IsExist(s.p.rootPath))
}

func (s *_S) TestGetLocalConfigureInfomation() {
	cacheData := configinfo.NewCacheData("linna", "DEFAULT_GROUP")
	content, err := s.p.GetLocalConfigureInformation(cacheData, false)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), content)
}

func (s *_S) TestGetLocalConfigureInfomationForce() {
	cacheData := configinfo.NewCacheData("linna", "DEFAULT_GROUP")
	content, err := s.p.GetLocalConfigureInformation(cacheData, true)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), content, "adafdasfdsa")
	println(content)
}

func (s *_S) TestGetFilePath() {
	println(s.p.getFilePath("linana", "DEFAULT_GROUP"))
}

func TestS(t *testing.T) {
	suite.Run(t, new(_S))
}
