package server

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	arr := strings.Split("my.test,DEFAULT_GROUP,", lineSeparator)
	fmt.Println(len(arr))
	fmt.Println(arr)

	builder := strings.Builder{}
	fmt.Println(builder.String())
}
