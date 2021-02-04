package stringutil

import (
	"crypto/md5"
	"fmt"
	"strings"
)

const invalidChar = ";&%#$@,*^~()/|\\+"

//HasInvalidChar whether str is empty or contain invalid char
func HasInvalidChar(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	}
	return strings.ContainsAny(str, invalidChar)
}

//GetMd5 get md5 from string
func GetMd5(content string) string {
	if content == "" {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
