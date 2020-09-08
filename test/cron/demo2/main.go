package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type job struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

// 使用1个调度协程，定时检查所有cron job，执行对应的job
func main() {
	jobs := make(map[string]*job)
	now := time.Now()
	expr1 := cronexpr.MustParse("*/2 * * * * * *")
	jobs["job1"] = &job{
		expr:     expr1,
		nextTime: expr1.Next(now),
	}

	expr2 := cronexpr.MustParse("*/5 * * * * * *")
	jobs["job2"] = &job{
		expr:     expr2,
		nextTime: expr2.Next(now),
	}

	done := make(chan bool)

	go func() {
		var counter int
		for {
			timer := time.NewTicker(time.Second)
			select {
			case <-timer.C:
				for jobName, job := range jobs {
					if job.nextTime.Unix() == time.Now().Unix() {
						fmt.Printf("%s is running\n", jobName)
						counter++
					}
				}
				if counter >= 2 {
					done <- true
				}
			}
		}
	}()

	select {
	case <-done:
		fmt.Println("done")
	}
}
