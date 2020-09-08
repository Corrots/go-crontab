package main

import (
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err    error
}

func main() {
	start := time.Now().Unix()
	resChan := make(chan result)
	done := make(chan bool)
	//cmd := exec.Command("/bin/bash", "-c", "sleep 2;echo hello")
	//var ctx context.Context
	//cancel, _ := context.WithCancel(ctx)
	command := "sleep 3; echo hello"

	go func() {
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		resChan <- result{
			output: output,
			err:    err,
		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("done")
	case <-time.After(time.Duration(start)):
		fmt.Println("cancel the command: ", command)
	case res := <-resChan:
		if res.err != nil {
			fmt.Println("err: ", res.err)
			return
		}
		fmt.Printf("output: %s\n", res.output)
	}
	//fmt.Printf("output: %s\n", output)
}
