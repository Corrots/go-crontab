package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// 分布式锁
// Lease实现锁自动过期
// op操作
// txn事务操作 if then else

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.56.23:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("new etcd client: %v\n", err)
		return
	}
	// Op Operation
	// 1.上锁(创建租约，自动续租，拿着租约去抢占一个key)
	// 2.业务处理
	// 3.释放锁(取消自动续租，释放租约)
	lease := clientv3.NewLease(client)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	// 创建租约
	leaseResp, err := lease.Grant(ctx, 10)
	if err != nil {
		fmt.Printf("lease grant err: %v\n", err)
		return
	}
	// 自动续租
	aliveChan, err := lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		fmt.Printf("lease keep alive err: %v\n", err)
		return
	}
	// 消费自动续租resp
	go func() {
	loop:
		for {
			select {
			case res := <-aliveChan:
				if res == nil {
					fmt.Printf("自动续租失败: %v\n", res.ID)
					break loop
				}
				fmt.Printf("自动续租成功: %v\n", res.ID)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	// 使用事务去抢占锁(key)
	key := "/cron/lock/lock3"
	//kv := clientv3.NewKV(client)
	txn := clientv3.NewKV(client).Txn(context.TODO())
	//cmp := clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "node:192.168.56.20", clientv3.WithLease(leaseResp.ID))).
		Else(clientv3.OpGet(key))
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Printf("txn commit err: %v\n", err)
		return
	}
	if !txnResp.Succeeded {
		fmt.Printf("抢占锁失败, [%s]占用中\n", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
		return
	}
	fmt.Println("抢占锁成功")
	fmt.Println("业务处理中")
	time.Sleep(time.Second * 100)

	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseResp.ID)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	select {
	case <-interrupt:
		signal.Stop(interrupt)
		lease.Revoke(context.TODO(), leaseResp.ID)
		cancelFunc()
	}
}
