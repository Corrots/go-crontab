package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/coreos/etcd/clientv3"
)

// watch
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
	go func() {
		for {
			kv.Put(context.TODO(), "/cronjob/job7", "I am job7")
			time.Sleep(time.Second)
			kv.Delete(context.TODO(), "/cronjob/job7")
		}
	}()

	// 获取开始监听的Revision
	getResp, err := kv.Get(context.TODO(), "/cronjob/job7")
	if err != nil {
		fmt.Printf("get by key err: %v\n", err)
		return
	}
	fmt.Printf("current value: %s\n", getResp.Kvs[0].Value)
	latestRevision := getResp.Header.Revision + 1
	fmt.Printf("start watch Revision: %d\n", latestRevision)
	ctx := context.TODO()
	var cancelFunc func()
	ctx, cancelFunc = context.WithCancel(ctx)
	time.AfterFunc(time.Second*5, func() {
		cancelFunc()
	})
	watcher := clientv3.NewWatcher(client)
	watchChan := watcher.Watch(ctx, "/cronjob/job7", clientv3.WithRev(latestRevision))

	go func() {
		for {
			select {
			case resp := <-watchChan:
				//fmt.Printf("watch resp: %+v\n", resp)
				for _, e := range resp.Events {
					switch e.Type {
					case mvccpb.PUT:
						fmt.Printf("PUT %s=%s\n", e.Kv.Key, e.Kv.Value)
					case mvccpb.DELETE:
						fmt.Printf("DELETE Revision: %d\n", e.Kv.ModRevision)
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.Sleep(time.Second * 120)
}
