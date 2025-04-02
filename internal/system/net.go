package system

import (
	"errors"
	"net"
)

func GetLocalIP() (string, error) {
	// 获取网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 如果接口未启用，跳过
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 如果接口是回环接口，跳过
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		// 遍历所有地址
		for _, addr := range addrs {
			var ip net.IP
			// 判断地址类型并获取IP地址
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 如果IP地址为空或是回环地址，跳过
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// 将IP地址转换为IPv4地址
			ip = ip.To4()
			// 如果转换后的IP地址为空，跳过
			if ip == nil {
				continue
			}

			// 返回找到的IP地址
			return ip.String(), nil
		}
	}

	// 如果未找到本地IP地址，返回错误
	return "", errors.New("找不到本地 IP 地址")
}
