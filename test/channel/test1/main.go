package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int, 10)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for {
			n, ok := <-ch
			if !ok {
				return
			}
			fmt.Println(n)
		}
	}()
	time.Sleep(time.Minute)
}
