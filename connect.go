package main

import (
	"fmt"
	"net"
	"time"
)

func Connect(ip string, port int) (string, int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port),
		time.Duration(timeout)*time.Second)
	if err != nil {
		return ip, port, err
	}
	defer conn.Close()

	return ip, port, err
}
