package processor

import (
	"fmt"
	"gdiamond/client/configinfo"
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
	os.Setenv("user.home", "./")
	c, err := configinfo.NewConfigure()
	assert.NoError(s.T(), err)
	p := NewServerAddressProcessor(c)
	s.p = p
	s.p.Start()
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

func (s *_SS) TestName() {
}
