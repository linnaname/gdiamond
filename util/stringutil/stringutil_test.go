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
