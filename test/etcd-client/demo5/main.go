package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// 分布式乐观锁相关
// 租约: Lease
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

	lease := clientv3.NewLease(client)
	ctx := context.TODO()
	grant, err := lease.Grant(ctx, 10)
	if err != nil {
		fmt.Printf("leases err: %v\n", err)
		return
	}
	// 通过context设置timeout在5S后取消自动续租，
	// 此时并不会马上终止租约，所以lease的租约总TTL=15s
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	// 为Lease续租
	// 自动续租
	aliveResp, err := lease.KeepAlive(ctx, grant.ID)
	if err != nil {
		fmt.Printf("lease KeepAlive err: %v\n", err)
		return
	}

	// 开启goroutine,消费续租response chan中的resp
	go func() {
	loop:
		for {
			select {
			case resp := <-aliveResp:
				if <-aliveResp == nil {
					fmt.Println("lease expired with nil keep alive resp chan")
					break loop
				}
				fmt.Println("keep alive succeed resp: ", resp)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	kv := clientv3.NewKV(client)
	putResp, err := kv.Put(ctx, "/cronjob/job5", "Dark Souls", clientv3.WithLease(grant.ID))
	if err != nil {
		fmt.Printf("put err: %v\n", err)
		return
	}
	//
	for {
		geResp, err := kv.Get(context.TODO(), "/cronjob/job5")
		if err != nil {
			fmt.Printf("get err: %v\n", err)
			return
		}
		if geResp.Count == 0 {
			fmt.Println("TTL reached")
			break
		}
		fmt.Printf("not expired: %v\n", geResp.Kvs)
		time.Sleep(time.Second)
	}

	fmt.Println("putResp: ", putResp)
}
