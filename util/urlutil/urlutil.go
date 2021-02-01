package urlutil

import (
	"fmt"
)

//GetURL concat to  url path
func GetURL(domainName string, port int, uri string) string {
	return fmt.Sprintf("http://%s:%v/%s", domainName, port, uri)
}
