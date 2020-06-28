package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/corrots/go-crontab/common"
)

const (
	EventPut    = 0
	EventDelete = 1
	EventKill   = 2
)

type Event struct {
	Type int
	Task *Task
}

// 获取当前所有task相关的event，并监听task的put/delete
func (e *Worker) CollectEvent() {
	getResp, err := e.KV.Get(context.TODO(), common.TaskNamePrefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("etcd get %s err: %v\n", common.TaskNamePrefix, err)
		return
	}
	if len(getResp.Kvs) == 0 {
		return
	}
	for _, v := range getResp.Kvs {
		var task Task
		if err := json.Unmarshal(v.Value, &task); err != nil {
			log.Printf("unmarshal task json err: %v\n", err)
			return
		}
		e.Scheduler.PushEvent(Event{Type: EventPut, Task: &task})
	}
	// 监听task的put/delete
	go func() {
		id := getResp.Header.Revision + 1
		watchChan := e.Watcher.Watch(context.TODO(), common.TaskNamePrefix, clientv3.WithPrefix(), clientv3.WithRev(id))
		for {
			select {
			case resp := <-watchChan:
				for _, v := range resp.Events {
					task := new(Task)
					if err := json.Unmarshal(v.Kv.Value, task); err != nil {
						log.Printf("unmarshal task json err: %v\n", err)
						return
					}
					fmt.Printf("%d %s=%s\n", v.Type, v.Kv.Key, v.Kv.Value)
					e.Scheduler.PushEvent(Event{Type: int(v.Type), Task: task})
				}
			}
		}
	}()
	// 监听task的kill
	go func() {
		watchChan := e.Watcher.Watch(context.TODO(), common.TaskKillerPrefix, clientv3.WithPrefix())
		for {
			select {
			case resp := <-watchChan:
				for _, v := range resp.Events {
					task := new(Task)
					if err := json.Unmarshal(v.Kv.Value, task); err != nil {
						log.Printf("unmarshal task json err: %v\n", err)
						return
					}
					if v.Type == clientv3.EventTypePut {
						fmt.Printf("Kill %s=%s\n", v.Kv.Key, v.Kv.Value)
						e.Scheduler.PushEvent(Event{Type: EventKill, Task: task})
					}
				}
			}
		}
	}()
}
