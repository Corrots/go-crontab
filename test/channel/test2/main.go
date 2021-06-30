package main

import (
	"fmt"
)

func main() {

	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				close(ch)
				return
			case n := <-ch:
				fmt.Println(n)
			}
		}
	}()
	done <- struct{}{}
}
