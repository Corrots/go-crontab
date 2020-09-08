package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	now := time.Now()
	parse, err := cronexpr.Parse(`*/2 * * * * * *`)
	if err != nil {
		panic(err)
	}
	next := parse.Next(now)
	fmt.Println(next.Format("2006-01-02 15:04:05"))

	time.AfterFunc(next.Sub(now), func() {
		fmt.Println("aaaa")
	})
	time.Sleep(time.Second * 5)
}
