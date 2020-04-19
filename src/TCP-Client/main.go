package main

import (
	"golang.org/x/exp/errors/fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8008")
	if err != nil {
		fmt.Println("Dial failed: ", err)
		return
	}
	defer conn.Close()

	go func() {
		buf := make([]byte, 4096)

		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println("Stdin read err: ", err)
				continue
			}
			conn.Write(buf[:n])
		}
	}()

	res := make([]byte, 4096)

	for {
		readCnt, err := conn.Read(res)
		if readCnt == 0 {
			fmt.Println("Server has closed...")
			return
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client has read done!")
				continue
			} else {
				fmt.Println("Client read err: ", err)
				return
			}
		}

		fmt.Println("Client read data: ", string(res[:readCnt]))
	}

}
