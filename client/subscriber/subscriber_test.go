package subscriber

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"gdiamond/util/maputil"
	"gdiamond/util/stringutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestCheckContent(t *testing.T) {
	assert.True(t, checkContent("test", stringutil.GetMd5("test")))
	assert.True(t, checkContent("", stringutil.GetMd5("")))
	assert.False(t, checkContent("A", stringutil.GetMd5("B")))
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

	h.Set(ContentEncoding, "gzip")
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte("gzip"))
	resp.Body = ioutil.NopCloser(strings.NewReader(b.String())) // r type is io.ReadCloser
	content = getContent(resp)
	assert.NotEmpty(t, content)
}

func TestIsZipContent(t *testing.T) {
	h := http.Header{}
	h.Set(ContentEncoding, "gzip")
	assert.True(t, isZipContent(h))
	h.Set(ContentEncoding, "zip")
	assert.False(t, isZipContent(h))
	h.Set(ContentEncoding, "")
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
	assert.Equal(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(0))
	s.s.AddDataId("linnana", "DEFAULT_GROUP")
	fmt.Println("AddDataId", s.s.cache)
	assert.Greater(s.T(), maputil.LengthOfSyncMap(s.s.cache), int64(0))
	value, loaded := s.s.cache.Load("linnana")
	assert.True(s.T(), loaded)
	assert.NotNil(s.T(), value)
	fmt.Println("value:", value)
	groupCache, ok := value.(sync.Map)
	assert.True(s.T(), ok)
	assert.NotNil(s.T(), groupCache)
	assert.Greater(s.T(), maputil.LengthOfSyncMap(groupCache), int64(0))
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

func (s *_S) TestGetConfigureInformation() {
	content := s.s.GetConfigureInformation("linna", "DEFAULT_GROUP", 1000)
	assert.NotEmpty(s.T(), content)
}

func (s *_S) TestGetAvailableConfigureInformation() {
	content := s.s.GetAvailableConfigureInformation("linna", "DEFAULT_GROUP", 1000)
	assert.NotEmpty(s.T(), content)
}

func (s *_S) TestGetAvailableConfigureInformationFromSnapshot() {
	content := s.s.GetAvailableConfigureInformationFromSnapshot("linna", "DEFAULT_GROUP", 1000)
	assert.NotEmpty(s.T(), content)
}
