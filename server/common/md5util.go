package common

import (
	"crypto/md5"
	"fmt"
)

func GetMd5(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
