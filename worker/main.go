package main

import (
	"log"
	"time"

	"github.com/corrots/go-crontab/worker/etcd"

	"github.com/corrots/go-crontab/common"
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

	client, err := etcd.NewWorker()
	if err != nil {
		log.Println(err)
		return
	}
	client.CollectEvent()
	client.Scheduler.Run()
	time.Sleep(time.Minute * 5)
}
