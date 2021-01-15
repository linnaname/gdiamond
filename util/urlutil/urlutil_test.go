package urlutil

import (
	"fmt"
	"testing"
)

func TestGetUrl(t *testing.T) {
	fmt.Println(GetUrl("www.linnana.me", 8081, "test"))
}
