package main

import (
	"log"
	"time"

	"github.com/corrots/go-crontab/common"
	"github.com/corrots/go-crontab/worker/manager"
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
	if err := manager.InitJobManager(); err != nil {
		log.Println(err)
		return
	}

	manager.InitScheduler()
	manager.JM.WatchJobs()
	manager.JM.WatchKiller()

	time.Sleep(time.Hour)
}
