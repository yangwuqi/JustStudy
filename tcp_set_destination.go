package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	var destination, markHost,listenPort string
	listenPort="10.15.165.4:808"

	if len(os.Args) < 2 {
		destination = "10.15.165.5:888"
		markHost = "10.15.165.4:808"
		fmt.Printf("default destination is 10.15.165.5:888, default markHost is 10.15.165.4:808, listenPort is %v\n",listenPort)
	} else if len(os.Args) == 3 {
		destination = os.Args[1]
		markHost = os.Args[2]
		fmt.Printf("destination is %v, markHost is %v, listenPort is %v\n", destination, markHost,listenPort)
	} else if len(os.Args) == 2 {
		destination = os.Args[1]
		markHost = "10.15.165.4:808"
		fmt.Printf("destination is %v, default markHost is 10.15.165.4:808, listenPort is %v\n", destination,listenPort)
	} else {
		fmt.Printf("your args is invalid!!!\n")
		return
	}

	fmt.Printf("开始！！！\n")
	listen, err := net.Listen("tcp", listenPort)
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
		go handle(conn, destination, markHost)
	}
}

func handle(conn net.Conn, destination string, markHost string) { //处理函数
	defer conn.Close()
	buf := make([]byte, 20000)
	//var buf []byte
	n, _ := conn.Read(buf)

	if n > 20000 {
		fmt.Println("获得的TCP流过大！缓冲区只有20000字节！")
		return
	}
	fmt.Printf("brower data: %v\n", string(buf))

	httpHostString := getHost(buf[:n])
	fmt.Printf("解析tcp流HTTP获得的host地址是 %v\n", httpHostString)

	if httpHostString == markHost {
		fmt.Println("该tcp流HTTP的host地址与被标记的host相同，开始设置跳转")
		jumpHandle(conn, buf, n, destination)
	} else {
		fmt.Printf("httpHostString和markHost不相等\n")
		return
	}

	//msg := []byte("the normal tcp write-back" + " --- " + string(buf[:n]))
	//conn.Write(msg) //正常TCP写回数据
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
	if len(host) == 0 {
		return string(host)
	} else {
		return string(host[:len(host)-1])
	}
}

func jumpHandle(conn net.Conn, buf []byte, n int, destination string) {
	conn2, err := net.Dial("tcp", destination)
	defer conn2.Close()
	if err != nil {
		fmt.Printf("跳转dialing错误！%v\n", err.Error())
		panic(err)
	}

	fmt.Printf("\njump and send: %v|||to destination %v\n", string(buf[:n]), destination)
	_, err2 := conn2.Write(buf[:n])
	if err2 != nil {
		fmt.Println("跳转时发送数据出错！")
		panic(err2)
	}
	buf2 := make([]byte, 20000)
	for {
		n2, err3 := conn2.Read(buf2)
		if err3 != nil {
				fmt.Println("跳转接收读数据读完啦！")
				break
		}

		fmt.Printf("\nsecond buf: %v\n", string(buf2))
		conn.Write(buf2[:n2])
	}
}
