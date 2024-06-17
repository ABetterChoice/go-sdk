package env

import (
	"net"
	"strconv"
	"strings"
)

// DefaultLocalIP Use this IP after failing to obtain the local IP
const DefaultLocalIP = "127.0.0.1"

var localIP string

func init() {
	localIP = DefaultLocalIP
	ips := getLocalIplist(1)
	if len(ips) > 0 {
		localIP = ips[0]
	}
}

// LocalIP local ip
func LocalIP() string {
	return localIP
}

func getLocalIplist(limit int) []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	ipList := getIpListFromAddr(limit, addrs)
	if len(ipList) == 0 {
		ipList = append(ipList, DefaultLocalIP)
	}
	return ipList
}

func getIpListFromAddr(limit int, addrs []net.Addr) []string {
	count := 0
	var ipList []string
	for _, address := range addrs {
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		if !isAddrOK(ipNet) {
			continue
		}
		ipList = append(ipList, ipNet.IP.String())
		count++
		if count >= limit {
			break
		}
	}
	return ipList
}

func isAddrOK(ipNet *net.IPNet) bool {
	if ipNet.IP.IsLoopback() {
		return false
	}
	if ipNet.IP.To4() != nil {
		if isInnerIp(ipNet.IP.String()) {
			return true
		}
	} else if ipNet.IP.To16() != nil {
		return true
	}
	return false
}

// isInnerIp Here we determine whether the URL is an intranet IP
func isInnerIp(ipv4 string) bool {
	temp := strings.Split(ipv4, ".")
	firstNum, _ := strconv.Atoi(temp[0])
	// 100 172 192 The beginning can be regarded as the intranet ip
	inValues := []int{100, 172, 192}
	for i := 0; i < len(inValues); i++ {
		if firstNum == inValues[i] {
			return true
		}
	}
	// The number starts from 1 to 15, and can also be an intranet IP.
	if firstNum >= 1 && firstNum <= 15 {
		return true
	}
	return false
}
