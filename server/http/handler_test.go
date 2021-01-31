package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	arr := strings.Split("my.test,DEFAULT_GROUP,", LINE_SEPARATOR)
	fmt.Println(len(arr))
	fmt.Println(arr)
}
