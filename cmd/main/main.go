package main

import (
	"fmt"
	"github.com/ark-go/yeelight/pkg/discovery"
	"net"
)

func init() {
}

func main() {
	fmt.Println("Запрос...")
	discovery.StartDiscover()
	fmt.Println(GetLocalIP())
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
