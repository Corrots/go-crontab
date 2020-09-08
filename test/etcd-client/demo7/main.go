package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// op 取代 get put delete
func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.56.23:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("new etcd client: %v\n", err)
		return
	}

	kv := clientv3.NewKV(client)
	// Op Operation
	key := "/cron/jobs/job8"
	op := clientv3.OpPut(key, "job8")
	opResp, err := kv.Do(context.TODO(), op)
	if err != nil {
		fmt.Printf("op err: %v\n", err)
		return
	}
	fmt.Printf("Put Revision: %d\n", opResp.Put().Header.Revision)

	op = clientv3.OpGet(key)
	opResp, err = kv.Do(context.TODO(), op)
	if err != nil {
		fmt.Printf("op err: %v\n", err)
		return
	}
	getRespKv := opResp.Get().Kvs[0]
	fmt.Printf(
		"Get %s=%s,Revision: %d\n",
		getRespKv.Key,
		getRespKv.Value,
		getRespKv.ModRevision)
	// OpDelete
	op = clientv3.OpDelete(key, clientv3.WithPrevKV())
	opResp, err = kv.Do(context.TODO(), op)
	if err != nil {
		fmt.Printf("op err: %v\n", err)
		return
	}
	delRespKvs := opResp.Del().PrevKvs[0]
	fmt.Printf(
		"Del %s=%s,Revision: %d\n",
		delRespKvs.Key,
		delRespKvs.Value,
		delRespKvs.ModRevision)
}
