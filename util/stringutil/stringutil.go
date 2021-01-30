package stringutil

import "strings"

const invalidChar = ";&%#$@,*^~()/|\\+"

//HasInvalidChar whether str is empty or contain invalid char
func HasInvalidChar(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	}
	return strings.ContainsAny(str, invalidChar)
}
