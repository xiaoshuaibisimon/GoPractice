package main

import (
	"golang.org/x/exp/errors/fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("listen err: ", err)
		return
	}
	defer listener.Close()

	for {
		fmt.Println("Server is waiting for client...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept err: ", err)
			continue
		}

		go tcpHandle(conn)
	}
}

func tcpHandle(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr()
	fmt.Println("Client addr is :", clientAddr)

	buf := make([]byte, 1024*4)

	for {
		readCnt, err := conn.Read(buf)
		if readCnt == 0 {
			fmt.Println("Client has closed...")
			return
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("Server has read done!")
				continue
			} else {
				fmt.Println("Server read err: ", err)
				return
			}
		}

		if "exit\n" == string(buf[:readCnt]) || "exit\r\n" == string(buf[:readCnt]) {
			fmt.Println("Client hope to close the connection.")
			return
		}

		conn.Write([]byte(strings.ToUpper(string(buf[:readCnt]))))
		fmt.Println("Server read data: ", string(buf[:readCnt]))
	}
}
