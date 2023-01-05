package utils

import (
	"fmt"
	"net"
)

// 获取可用端口
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "127.0.0.1"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// 判断端口是否可以（未被占用）
func IsAddressAvailable(address string) bool {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// log.Infof("port %s is taken: %s", address, err)
		return false
	}
	defer listener.Close()
	return true
}
