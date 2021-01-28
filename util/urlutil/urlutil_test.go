package urlutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUrl(t *testing.T) {
	assert.Equal(t, GetUrl("www.linnana.me", 8081, "test"), "http://www.linnana.me:8081/test")
}
