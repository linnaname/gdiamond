package stringutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasInvalidChar(t *testing.T) {
	assert.False(t, HasInvalidChar("abacd34343jkjkda--"))
	assert.True(t, HasInvalidChar("@dajkdjakdjka"))
	assert.True(t, HasInvalidChar("dafdsad\\jdajkldjak"))
}

func TestGetMd5(t *testing.T) {
	assert.Empty(t, GetMd5(""))
	assert.NotEmpty(t, GetMd5("abac"))
}
