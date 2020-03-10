package main

import (
	"golang.org/x/exp/errors/fmt"
	"net"
	"os"
)

func recvFile(conn net.Conn, fileName string) {
	//创建文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("os.Create err :", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("接收文件完毕")
			return
		}
		if err != nil {
			fmt.Println("os.Create err :", err)
			return
		}
		f.Write(buf[:n])
	}
}

func main() {
	//1.创建用于监听的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}
	defer listener.Close()

	//2.阻塞监听
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("listener.Accept,err:", err)
		return
	}
	defer conn.Close()

	//3.获取文件名
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read,err:", err)
		return
	}
	fileName := string(buf[:n])

	//4.回写OK给客户端
	conn.Write([]byte("OK"))

	//5.接收文件内容
	recvFile(conn, fileName)
}
