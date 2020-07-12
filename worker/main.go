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

	worker, err := core.NewWorker()
	if err != nil {
		log.Println(err)
		return
	}
	worker.CollectEvent()
	worker.Scheduler.Run()
	time.Sleep(time.Minute * 5)
}
