package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := ParseCmdAndInitConfig()
	assert.NoError(t, err)
	assert.NotNil(t, GMySQLConfig)
	assert.NotEmpty(t, GMySQLConfig.DBUrl)
	assert.Equal(t, GMySQLConfig.MaxIdleConns, 20)
	assert.Equal(t, GMySQLConfig.MaxOpenConns, 20)
}
