package urlutil

import (
	"fmt"
)

/**
GetUrl concat to  url path
*/
func GetUrl(domainName string, port int, uri string) string {
	return fmt.Sprintf("http://%s:%v/%s", domainName, port, uri)
}
