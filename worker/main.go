package main

import (
	"log"
	"time"

	"github.com/corrots/go-crontab/common"
	"github.com/corrots/go-crontab/worker/core"
	flag "github.com/spf13/pflag"
)

var (
	conf = flag.String("c", "../conf.yaml", "specific the service config file")
)

func main() {
	flag.Parse()
	if err := common.InitConfig(conf); err != nil {
		log.Println(err)
		return
	}
	// 服务注册和发现
	register, err := core.NewRegister()
	if err != nil {
		log.Println(err)
		return
	}
	if err := register.KeepOnline(); err != nil {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				register.KeepOnline()
			}
		}
	}
	//
	worker, err := core.NewWorker()
	if err != nil {
		log.Println(err)
		return
	}

	worker.CollectEvent()
	worker.Scheduler.Run()
	//time.Sleep(time.Minute * 5)
}
