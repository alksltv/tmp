package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

func TCPClient(a string, p string) {
	conn, err := net.DialTimeout("tcp", a+":"+p, 5*time.Second)
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		fmt.Println("TCP connection refused from: ", a+":"+p)
		return
	} else if err != nil && strings.Contains(err.Error(), "timeout") {
		fmt.Println("TCP connection timed out: ", a+":"+p)
		return
	} else if err != nil {
		fmt.Println("TCP connection failed: ", a+":"+p)
	} else {
		fmt.Println("TCP connected to: ", a+":"+p)
	}
	conn.Close()
}

func main() {
	clientIP := flag.String("ip", "127.0.0.1", "example: -ip=127.0.0.1")
	clientPort := flag.String("port", "8080", "example: -port=8080")
	flag.Parse()
	defer TCPClient(*clientIP, *clientPort)
}
