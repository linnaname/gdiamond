package service

import (
	"fmt"
	"gdiamond/server/internal/common"
	"gdiamond/server/internal/model"
	"gdiamond/util/maputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type _SS struct {
	suite.Suite
}

func (s *_SS) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	common.ParseCmdAndInitConfig()
	common.InitDBConn()
}

func (s *_SS) TearDownSuite() {
	common.CloseConn()
}

func TestSS(t *testing.T) {
	suite.Run(t, new(_SS))
}

func (s *_SS) TestAddConfigInfo() {
	err := AddConfigInfo("gdiamond.test.vv", "GDIADMOND", "whether true")
	assert.NoError(s.T(), err)
}

func (s *_SS) TestUpdateConfigInfo() {
	err := UpdateConfigInfo("gdiamond.test.vv", "GDIADMOND", "whether update ok")
	assert.NoError(s.T(), err)
}

func (s *_SS) TestFindConfigInfo() {
	cInfo, err := FindConfigInfo("gdiamond.test.vv", "GDIADMOND")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), cInfo)
	assert.NotEmpty(s.T(), cInfo.Content)
}

func (s *_SS) TestFindConfigInfoPage() {
	page, err := FindConfigInfoPage(1, 10, "GDIADMOND", "gdiamond.test.vv")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.Equal(s.T(), page.TotalCount, 1)

	pEmptyDataID, err := FindConfigInfoPage(1, 10, "GDIADMOND", "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyDataID)
	assert.Greater(s.T(), pEmptyDataID.TotalCount, 0)

	pEmptyGroup, err := FindConfigInfoPage(1, 10, "", "gdiamond.test.vv")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyGroup)
	assert.Greater(s.T(), pEmptyGroup.TotalCount, 0)

	pEmptyGroupAndDataID, err := FindConfigInfoPage(1, 10, "", "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), pEmptyGroupAndDataID)
	assert.Greater(s.T(), pEmptyGroupAndDataID.TotalCount, 0)
}

func (s *_SS) TestFindConfigInfoLike() {
	page, err := FindConfigInfoLike(1, 2, "lin", "DEFAULT_GROUP")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page)
	assert.NotNil(s.T(), page.PageItems)
	assert.Greater(s.T(), page.TotalCount, 1)
}

func TestGetConfigInfoPath(t *testing.T) {
	filePath := GetConfigInfoPath("linna", "DEFAULT_GROUP")
	assert.NotEmpty(t, filePath)
	assert.Contains(t, filePath, "linna")
}

func TestUpdateMD5Cache(t *testing.T) {
	assert.Equal(t, maputil.LengthOfSyncMap(cache), int64(0))
	cInfo := model.NewConfigInfo("linna", "DEFAULT_GROUP", "gdiamond", time.Now())
	UpdateMD5Cache(cInfo)
	assert.Equal(t, maputil.LengthOfSyncMap(cache), int64(1))
}

func TestGetContentMD5(t *testing.T) {
	assert.Equal(t, maputil.LengthOfSyncMap(cache), int64(0))
	cInfo := model.NewConfigInfo("linna", "DEFAULT_GROUP", "gdiamond", time.Now())
	UpdateMD5Cache(cInfo)
	md5 := GetContentMD5("linna", "DEFAULT_GROUP")
	assert.NotEmpty(t, md5)
}
