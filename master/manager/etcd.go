package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
)

var (
	JM *JobManager
)

type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func InitJobManager() error {
	config := clientv3.Config{
		Endpoints:   strings.Split(viper.GetString("etcd.endpoints"), ";"),
		DialTimeout: time.Duration(viper.GetInt64("etcd.dialTimeout")) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return fmt.Errorf("new etcd client err: %v\n", err)
	}

	JM = &JobManager{
		client: client,
		kv:     clientv3.NewKV(client),
		lease:  clientv3.NewLease(client),
	}
	return nil
}
