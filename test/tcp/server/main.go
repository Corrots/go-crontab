package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			for {
				var buf [128]byte
				_, err := c.Read(buf[:])
				if err != nil {
					fmt.Printf("read from conn err: %v\n", err)
					break
				}
				str := string(buf[:])
				fmt.Printf("From client, data: %v\n", str)
				_, err = c.Write(buf[:])
				if err != nil {
					fmt.Printf("write to client err: %v\n", err)
					break
				}
			}
		}(conn)
	}
}
