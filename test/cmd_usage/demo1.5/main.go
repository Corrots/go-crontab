package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err    error
}

func main() {
	start := time.Now()
	resChan := make(chan result)
	done := make(chan bool)
	command := "sleep 3; echo hello"

	ctx, cancelFunc := context.WithCancel(context.TODO())

	go func() {
		cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
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
	case <-time.After(time.Second * 1):
		fmt.Println("since: ", time.Since(start).Milliseconds())
		fmt.Println("cancel the command: ", command)
		cancelFunc()
	case res := <-resChan:
		if res.err != nil {
			fmt.Println("err: ", res.err)
			return
		}
		fmt.Printf("output: %s\n", res.output)
	}
	//fmt.Printf("output: %s\n", output)
}
