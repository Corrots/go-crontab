package core

import (
	"fmt"
	"net"

	"github.com/spf13/viper"
)

// 获取 WorkerID, 可获取到IP则使用IP作为workerId，否则使用配置文件中 serverIP
func getWorkerIP() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			// addr is ipv4/ipv6/unix socket
			ipNet, ok := addr.(*net.IPNet)
			// 是IP地址，且不是环回地址
			if ok && !ipNet.IP.IsLoopback() {
				// 只需要IPv4地址
				ipv4 := ipNet.IP.To4()
				if ipv4 != nil {
					fmt.Printf("local IP: %s\n", ipv4.String())
					return ipv4.String()
				}
			}
		}
	}
	return viper.GetString("worker:serverIP")
}
