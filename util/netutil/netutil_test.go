package netutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	assert.NotEmpty(t, ip)
	fmt.Println(ip)
}
