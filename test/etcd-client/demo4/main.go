package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// create etcd client and Delete key value
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
	deleteResponse, err := kv.Delete(context.TODO(), "/cronjob/job1", clientv3.WithPrefix(), clientv3.WithPrevKV())
	if err != nil {
		fmt.Printf("delete err: %v\n", err)
		return
	}
	//
	fmt.Println("PrevKvs: ", deleteResponse.PrevKvs)
	if len(deleteResponse.PrevKvs) != 0 {
		fmt.Printf("delete %+v\n", deleteResponse.PrevKvs)
	}
	fmt.Println("deleteResponse: ", deleteResponse)
}
