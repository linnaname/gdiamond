package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitDBConn(t *testing.T) {
	ParseCmdAndInitConfig()
	err := InitDBConn()
	assert.NoError(t, err)
	assert.NotNil(t, GDBConn)
	CloseConn()
}
