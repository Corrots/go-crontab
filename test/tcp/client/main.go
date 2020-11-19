package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 读取标准输入
	input := bufio.NewReader(os.Stdin)
	for {
		s, err := input.ReadString('\n')
		if err != nil {
			fmt.Printf("read from terminal err: %v\n", err)
			break
		}
		// 输入Q即退出
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			break
		}
		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Printf("write from conn err: %v\n", err)
			break
		}
		var buf [1024]byte
		_, err = conn.Read(buf[:])
		if err != nil {
			fmt.Printf("Read from conn err: %v\n", err)
			break
		}
		io.WriteString(os.Stdout, fmt.Sprintf("From server: %s\n", buf))
	}
}
