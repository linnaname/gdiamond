package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMd5(t *testing.T) {
	assert.Empty(t, GetMd5(""))
	assert.NotEmpty(t, GetMd5("abac"))
}
