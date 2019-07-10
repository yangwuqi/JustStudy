package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Printf("开始！！！\n")

	listen, err := net.Listen("tcp", "10.15.165.4:8080")

	if err != nil {
		fmt.Printf("监听出错！\n")
		panic(err)
	}

	for {
		conn, err1 := listen.Accept()
		if err1 != nil {
			fmt.Printf("监听accept出错！\n")
			continue
		}
		buf := make([]byte, 100000)
		n, _ := conn.Read(buf)
		fmt.Printf("从这个host地址收到了tcp流 %v:\n%v", getHost(buf), string(buf[:n]))
	}
}

func getHost(buf []byte) string {
	var host []byte
	for i := 0; i < len(buf)-6; i++ {
		if buf[i] == 'H' && buf[i+1] == 'o' &&
			buf[i+2] == 's' && buf[i+3] == 't' &&
			buf[i+4] == ':' && buf[i+5] == ' ' {
			for j := i + 6; buf[j] != '\n' && buf[j] != ' '; j++ {
				host = append(host, buf[j])
			}
			break
		}
	}
	return string(host[:len(host)-1])
}
