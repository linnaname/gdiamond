package configinfo

import (
	"gdiamond/util/fileutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfigure(t *testing.T) {
	os.Setenv("user.home", "./")
	c, err := NewConfigure()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.True(t, fileutil.IsExist(c.GetFilePath()))
}

func TestConfigure_AddDomainName(t *testing.T) {
	os.Setenv("user.home", "./")
	c, err := NewConfigure()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	c.AddDomainName("test")
	assert.NotEmpty(t, c.GetDomainNameList())
}
