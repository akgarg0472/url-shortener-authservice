package utils

import (
	"net"
	"os"
)

func GetHostIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()

		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

func GetHostName() string {
	hostname, err := os.Hostname()

	if err != nil {
		return "localhost"
	}

	return hostname
}
