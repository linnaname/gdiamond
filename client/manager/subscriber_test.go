package manager

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gdiamond/server/common"
	"gdiamond/util/maputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCheckContent(t *testing.T) {
	assert.True(t, checkContent("test", common.GetMd5("test")))
	assert.True(t, checkContent("", common.GetMd5("")))
	assert.False(t, checkContent("A", common.GetMd5("B")))
}

func TestGetContent(t *testing.T) {
	resp := &http.Response{}
	h := http.Header{}
	resp.Header = h
	resp.Body = ioutil.NopCloser(strings.NewReader("hello world")) // r type is io.ReadCloser
	content := getContent(resp)
	assert.NotEmpty(t, content)

	resp.Body = ioutil.NopCloser(strings.NewReader("")) // r type is io.ReadCloser
	content = getContent(resp)
	assert.Empty(t, content)

	h.Set(CONTENT_ENCODING, "gzip")
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte("gzip"))
	resp.Body = ioutil.NopCloser(strings.NewReader(b.String())) // r type is io.ReadCloser
	content = getContent(resp)
	assert.NotEmpty(t, content)
}

func TestIsZipContent(t *testing.T) {
	h := http.Header{}
	h.Set(CONTENT_ENCODING, "gzip")
	assert.True(t, isZipContent(h))
	h.Set(CONTENT_ENCODING, "zip")
	assert.False(t, isZipContent(h))
	h.Set(CONTENT_ENCODING, "")
	assert.False(t, isZipContent(h))
}

func TestConvertStringToSet(t *testing.T) {
	modifiedDataIdsString := ""
	set := convertStringToSet(modifiedDataIdsString)
	assert.Nil(t, set)
	modifiedDataIdsString = "OK"
	set = convertStringToSet(modifiedDataIdsString)
	assert.Equal(t, set.Size(), 0)

}

type _S struct {
	suite.Suite
	s *Subscriber
}

func (s *_S) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	s.s = GetSubscriberInstance()
	s.s.Start()
}

func (s *_S) TearDownSuite() {
	s.s.Close()
	fmt.Printf("TearDownSuite() ...\n")
}

func TestS(t *testing.T) {
	suite.Run(t, new(_S))
}

func (s *_S) TestGetSubscriberListener() {
	assert.NotNil(s.T(), s.s.GetSubscriberListener())
}

func (s *_S) TestAddDataId() {
	s.s.AddDataId("linnana", "DEFAULT_GROUP")
	fmt.Println("AddDataId", s.s.cache)
	assert.Greater(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(0))
}

func (s *_S) TestGetDataIds() {
	s.s.AddDataId("linnana", "DEFAULT_GROUP")
	s.s.AddDataId("me", "DEFAULT_GROUP")
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(2))
	dataIds := s.s.GetDataIds()
	assert.Equal(s.T(), dataIds.Size(), 2)
}

func (s *_S) TestRemoveDataId() {
	s.s.AddDataId("linnana", "DEFAULT_GROUP")
	s.s.AddDataId("me", "DEFAULT_GROUP")
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(2))
	s.s.RemoveDataId("linnana", "DEFAULT_GROUP")
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(1))

	s.s.RemoveDataId("test", "test")
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(1))
}

func (s *_S) TestGetConfigureInfomation() {
	content := s.s.GetConfigureInfomation("linnana", "DEFAULT_GROUP", 1000)
	fmt.Println("GetConfigureInfomation", content)
}

func (s *_S) TestGetAvailableConfigureInfomation() {
	content := s.s.GetAvailableConfigureInfomation("linnana", "DEFAULT_GROUP", 1000)
	fmt.Println("GetAvailableConfigureInfomation", content)
}

func (s *_S) TestGetAvailableConfigureInfomationFromSnapshot() {
	content := s.s.GetAvailableConfigureInfomationFromSnapshot("linnana", "DEFAULT_GROUP", 1000)
	fmt.Println("GetAvailableConfigureInfomationFromSnapshot", content)
}
