package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// create etcd client and Get key value with Prefix
func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.56.23:2379"},
		DialTimeout: time.Second * 5,
	})
	if err == context.DeadlineExceeded {
		fmt.Printf("DeadlineExceeded: %v\n", err)
		return
	}
	if err != nil {
		fmt.Printf("new etcd client: %v\n", err)
		return
	}

	kv := clientv3.NewKV(client)
	getResp, err := kv.Get(context.TODO(), "/cronjob/", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("put kv err: %v\n", err)
		return
	}
	if getResp.Count == 0 {
		fmt.Println("count is ", getResp.Count)
		return
	}
	fmt.Printf("get resp kvs: %+v\n\n", getResp.Kvs)
	for _, val := range getResp.Kvs {
		fmt.Printf("val: %s\n", val)
	}
}
