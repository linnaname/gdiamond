package subscriber

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSubscriberInstance(t *testing.T) {
	s := GetSubscriberInstance()
	assert.NotNil(t, s)
	s1 := GetSubscriberInstance()
	assert.Equal(t, s1, s)
}
