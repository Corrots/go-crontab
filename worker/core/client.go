package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

type Worker struct {
	Client    *clientv3.Client
	KV        clientv3.KV
	Lease     clientv3.Lease
	Watcher   clientv3.Watcher
	Scheduler Scheduler
}

func newEtcd() (*clientv3.Client, error) {
	config := clientv3.Config{
		Endpoints:   strings.Split(viper.GetString("etcd.endpoints"), ";"),
		DialTimeout: time.Duration(viper.GetInt64("etcd.dialTimeout")) * time.Millisecond,
	}
	return clientv3.New(config)
}

func NewWorker() (*Worker, error) {
	c, err := newEtcd()
	if err != nil {
		return nil, fmt.Errorf("etcd init err: %v\n", err)
	}
	return &Worker{
		Client:    c,
		KV:        clientv3.NewKV(c),
		Lease:     clientv3.NewLease(c),
		Watcher:   clientv3.NewWatcher(c),
		Scheduler: NewScheduler(),
	}, nil
}
