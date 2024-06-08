package utils

import (
	"fmt"
	"log"
	"net"
)

// 判断是否为私有 IP 地址
func isPrivateIP(ip net.IP) bool {
	// 私有地址段范围
	privateIPBlocks := []*net.IPNet{
		parseCIDR("10.0.0.0/8"),
		parseCIDR("172.16.0.0/12"),
		parseCIDR("192.168.0.0/16"),
		parseCIDR("100.64.0.0/10"),
		parseCIDR("192.0.0.0/24"),
		parseCIDR("192.0.2.0/24"),
		parseCIDR("192.88.99.0/24"),
		parseCIDR("198.18.0.0/15"),
		parseCIDR("203.0.113.0/24"),
		parseCIDR("240.0.0.0/4"),
	}

	// 检查 IP 是否在私有地址段范围内
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// 解析 CIDR 地址
func parseCIDR(cidr string) *net.IPNet {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicln(err)
	}
	return ipNet
}

func Ips() []string {
	var ips []string
	// 获取所有网络接口的信息
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// 遍历每个网络接口
	for _, faces := range interfaces {
		// 排除 loopback 接口
		if faces.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取该接口的地址信息
		addrs, err := faces.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// 遍历每个地址
		for _, addr := range addrs {
			// 检查地址类型是否是 IP 地址
			if ipNet, ok := addr.(*net.IPNet); ok {
				// 检查地址是否是 IPv4，并且是私有地址
				if ipNet.IP.To4() != nil && isPrivateIP(ipNet.IP) {
					ips = append(ips, ipNet.IP.String())
				}
			}
		}
	}
	ips = append(ips, "127.0.0.1")
	return ips
}
