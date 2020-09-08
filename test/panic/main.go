package main

import (
	"fmt"
	"runtime"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Println(string(buf))
		}
	}()
	panic("error occurred")
}
