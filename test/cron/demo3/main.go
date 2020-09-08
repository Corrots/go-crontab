package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type job struct {
	name     string
	expr     *cronexpr.Expression
	nextTime time.Time
}

// 使用1个调度协程，定时检查所有cron job，执行对应的job
func main() {
	var jobs []*job
	now := time.Now()
	expr1 := cronexpr.MustParse("*/2 * * * * * *")
	job1 := &job{
		name:     "job1",
		expr:     expr1,
		nextTime: expr1.Next(now),
	}

	expr2 := cronexpr.MustParse("*/5 * * * * * *")
	job2 := &job{
		name:     "job2",
		expr:     expr2,
		nextTime: expr2.Next(now),
	}
	jobs = append(jobs, job1, job2)

	go func() {
		for {
			now := time.Now()
			timer := time.NewTicker(time.Second)
			select {
			case <-timer.C:
				for _, job := range jobs {
					if job.nextTime.Equal(now) || job.nextTime.Before(now) {
						// 创建1个协程来执行job
						go func(name string) {
							fmt.Printf("%s executed\n", name)
							// 计算下一次调度时间
							job.nextTime = job.expr.Next(now)
							fmt.Println("next time: ", job.nextTime.Format("2006-01-02 15:04:05"))
						}(job.name)
						//
					}
				}
			}
			select {
			case <-time.NewTimer(time.Millisecond * 100).C:
			}
		}
	}()

	time.Sleep(time.Second * 15)
}
