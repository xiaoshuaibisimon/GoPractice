package main

import (
	"golang.org/x/exp/errors/fmt"
	"io"
	"net"
	"os"
)

func sendFile(conn net.Conn, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("os.Open err:", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 4096)

	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("发送文件完成")
			} else {
				fmt.Println("os.Read err:", err)
			}
			return
		}

		conn.Write(buf[:n])
	}

}

func main() {
	//1.获取文件名
	cmdList := os.Args
	if len(cmdList) != 2 {
		fmt.Println("go run xxx.go filepath")
		return
	}

	filePath := cmdList[1]

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("os.Stat,err:", err)
		return
	}

	fileName := fileInfo.Name()

	//2.主动发起连接
	conn, err := net.Dial("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("net.Dial,err:", err)
		return
	}
	defer conn.Close()

	//3.发送文件名给服务器
	conn.Write([]byte(fileName))

	//4.读取服务器返回的OK
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read,err:", err)
		return
	}

	if "OK" == string(buf[:n]) {
		//5.发文件内容给服务器
		sendFile(conn, filePath)
	} else {
		fmt.Println("bad response, ", string(buf[:n]))
	}

}
