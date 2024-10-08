package common

import (
	"fmt"
	"net"
)

func GetLocalIPv4() (string, error) {
	// 获取所有网络接口的地址列表
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %v", err)
	}

	// 遍历地址，寻找第一个非回环的IPv4地址
	for _, addr := range addrs {
		// 转换为IP地址类型
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			// 返回IPv4地址
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 address found")
}
