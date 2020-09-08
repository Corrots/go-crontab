package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// create etcd client and Put key value
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
	fmt.Println("succeed")
	ctx := context.TODO()
	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(ctx, "/cronjob/job4", "cyberpunk", clientv3.WithPrevKV())
	if err != nil {
		fmt.Printf("put kv err: %v\n", err)
		return
	}
	fmt.Println("put succeed Revision: ", putResp.Header.Revision)
	if putResp.PrevKv != nil {
		fmt.Println("put PrevKv: ", putResp.PrevKv)
	}
}
