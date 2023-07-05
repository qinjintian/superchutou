package net

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// GetExtranetIP 获取本机外网IP
func GetExtranetIP() (string, error) {
	resp, err := http.Get("http://ifconfig.me") // 获取外网 IP
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return fmt.Sprintf("%s", string(body)), nil
}

// GetLocalIP 获取本机内网IP
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if IPNet, ok := addr.(*net.IPNet); ok && !IPNet.IP.IsLoopback() {
			if IPNet.IP.To4() != nil {
				return IPNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("unable to get local ip")
}
