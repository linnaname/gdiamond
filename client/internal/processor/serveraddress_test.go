package processor

import (
	"fmt"
	"gdiamond/client/internal/configinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type _SS struct {
	suite.Suite
	p *ServerAddressProcessor
}

func (s *_SS) SetupSuite() {
	fmt.Printf("SetupSuite() ...\n")
	os.Setenv("user.home", "/Users/goranka")
	c, err := configinfo.NewConfigure()
	assert.NoError(s.T(), err)
	p := NewServerAddressProcessor(c)
	s.p = p
}

func (s *_SS) TearDownSuite() {
	fmt.Printf("TearDownSuite() ...\n")
	s.p.Stop()
}

func TestSS(t *testing.T) {
	suite.Run(t, new(_SS))
}

func TestGenerateLocalFilePath(t *testing.T) {
	path := generateLocalFilePath("test", "name")
	assert.NotEmpty(t, path)
	println(path)
}

func (s *_SS) TestAcquireServerAddressOnce() {
	p := s.p
	assert.True(s.T(), p.diamondConfigure.GetDomainNameList().Empty())
	b := p.acquireServerAddressOnce(0)
	assert.True(s.T(), b)
	assert.False(s.T(), p.diamondConfigure.GetDomainNameList().Empty())
	assert.Equal(s.T(), p.diamondConfigure.GetDomainNameList().Size(), 1)
	ele, ok := p.diamondConfigure.GetDomainNameList().Get(0)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), ele, "127.0.0.1")
}

func (s *_SS) TestStoreServerAddressesToLocal() {
	p := s.p
	b := p.acquireServerAddressOnce(0)
	assert.True(s.T(), b)
	p.storeServerAddressesToLocal()
}

func (s *_SS) TestReloadServerAddresses() {
	p := s.p
	assert.True(s.T(), p.diamondConfigure.GetDomainNameList().Empty())
	p.reloadServerAddresses()
	assert.False(s.T(), p.diamondConfigure.GetDomainNameList().Empty())
}

func (s *_SS) TestAcquireServerAddressFromLocal() {
	p := s.p
	p.acquireServerAddressFromLocal()
}
