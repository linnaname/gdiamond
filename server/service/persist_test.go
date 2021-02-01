package service

import (
	"crypto/md5"
	"fmt"
	"gdiamond/server/common"
	"gdiamond/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type _Suite struct {
	suite.Suite
}

func (s *_Suite) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	common.InitConfig()
	common.InitDBConn()
}

func (s *_Suite) TearDownSuite() {
	common.CloseConn()
}

func (s *_Suite) TestAddConfingInfo() {
	configInfo := &model.ConfigInfo{Group: "AAA_GROUP", DataID: "linna.com", Content: "song for linana", MD5: fmt.Sprintf("%x", md5.Sum([]byte("song for linana")))}
	err := addConfigInfo(configInfo)
	assert.NoError(s.T(), err)
}

func (s *_Suite) TestUpdateConfigInfo() {
	configInfo := &model.ConfigInfo{Group: "DEFAULT_GROUP", DataID: "linname", Content: "I hate linnana too", MD5: fmt.Sprintf("%x", md5.Sum([]byte("I hate linnana too")))}
	err := updateConfigInfo(configInfo)
	assert.NoError(s.T(), err)
}

func (s *_Suite) TestFindConfigInfo() {
	configInfo, err := findConfigInfo("linname", "DEFAULT_GROUP")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), configInfo)
	assert.Equal(s.T(), configInfo.Content, "what happend")
	assert.Equal(s.T(), configInfo.Group, "DEFAULT_GROUP")
	assert.Equal(s.T(), configInfo.MD5, "8d089c6892a42c2d1786f40ed8063850")
	assert.NotNil(s.T(), configInfo.LastModified)
}

func (s *_Suite) TestFindConfigInfoById() {
	configInfo, err := findConfigInfoByID(3)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), configInfo)
	assert.Equal(s.T(), configInfo.DataID, "linna")
	assert.Equal(s.T(), configInfo.ID, int64(3))
	assert.Equal(s.T(), configInfo.Content, "adafdasfdsa")
	assert.Equal(s.T(), configInfo.Group, "DEFAULT_GROUP")
	assert.Equal(s.T(), configInfo.MD5, "b9e1dbe39c2fd1e9c573b20de190c5bf")
	assert.NotNil(s.T(), configInfo.LastModified)
}

func (s *_Suite) TestFindConfigInfoByDataId() {
	page, err := findConfigInfoByDataID(1, 1, "linna.com")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 2)
	assert.Equal(s.T(), page.PageAvailable, 2)
	assert.Equal(s.T(), page.PageNO, 1)
	assert.Len(s.T(), page.PageItems, 1)
	configInfo, ok := page.PageItems[0].(*model.ConfigInfo)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), configInfo.Content, "song for linana")
}

func (s *_Suite) TestFindAllConfigInfo() {
	page, err := findAllConfigInfo(1, 10)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 4)
	assert.Equal(s.T(), page.PageAvailable, 1)
	assert.Equal(s.T(), page.PageNO, 1)
	assert.Len(s.T(), page.PageItems, 4)
	configInfo, ok := page.PageItems[0].(*model.ConfigInfo)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), configInfo.Content, "I hate linnana too")
}

func (s *_Suite) TestFindAllConfigLikeGroup() {
	page, err := findAllConfigLike(1, 10, "", "DEFAULT")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 3)
	assert.Equal(s.T(), page.PageAvailable, 1)
	assert.Equal(s.T(), page.PageNO, 1)
	assert.Len(s.T(), page.PageItems, 3)
}

func (s *_Suite) TestFindAllConfigLikeDataId() {
	page, err := findAllConfigLike(1, 10, "com", "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 2)
	assert.Equal(s.T(), page.PageAvailable, 1)
	assert.Equal(s.T(), page.PageNO, 1)
	assert.Len(s.T(), page.PageItems, 2)
}

func (s *_Suite) TestFindAllConfigLikeGroupAndDataId() {
	page, err := findAllConfigLike(1, 10, "com", "AAA")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 1)
	assert.Equal(s.T(), page.PageAvailable, 1)
	assert.Equal(s.T(), page.PageNO, 1)
	assert.Len(s.T(), page.PageItems, page.TotalCount)
}

func (s *_Suite) TestRemoveConfigInfo() {
	configInfo := &model.ConfigInfo{Group: "AAA_GROUP", DataID: "linna.com", Content: "song for linana", MD5: fmt.Sprintf("%x", md5.Sum([]byte("song for linana")))}
	err := removeConfigInfo(configInfo)
	assert.NoError(s.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(_Suite))
}
