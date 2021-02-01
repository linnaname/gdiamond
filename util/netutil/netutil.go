package netutil

import (
	"net"
)

//GetLocalIP get local ip, if more than one return first one, if can't get it return empty string
func GetLocalIP() string {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range interfaceAddrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
