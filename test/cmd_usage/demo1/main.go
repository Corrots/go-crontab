package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c","sleep 2;echo hello")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("output: %s\n", output)
}
