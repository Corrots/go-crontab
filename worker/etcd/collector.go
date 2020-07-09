package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
func (w *Worker) CollectEvent() {
	getResp, err := w.KV.Get(context.TODO(), common.TaskNamePrefix, clientv3.WithPrefix())
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
		w.Scheduler.PushEvent(Event{Type: EventPut, Task: &task})
	}
	// 监听task的put/delete
	go func(w *Worker) {
		id := getResp.Header.Revision + 1
		watchChan := w.Watcher.Watch(context.TODO(), common.TaskNamePrefix, clientv3.WithPrefix(), clientv3.WithRev(id))
		for {
			select {
			case resp := <-watchChan:
				for _, v := range resp.Events {
					var e Event
					var task Task
					switch v.Type {
					case EventPut:
						if err := json.Unmarshal(v.Kv.Value, &task); err != nil {
							log.Printf("unmarshal task json err: %v\n", err)
							continue
						}
						e = Event{Type: EventPut, Task: &task}
					case EventDelete:
						name := strings.TrimPrefix(string(v.Kv.Key), common.TaskNamePrefix)
						e = Event{Type: EventDelete, Task: &Task{Name: name}}
					}
					//fmt.Printf("%d %s=%s\n", v.Type, v.Kv.Key, v.Kv.Value)
					w.Scheduler.PushEvent(e)
				}
			}
		}
	}(w)
	// 监听task的kill
	go func() {
		watchChan := w.Watcher.Watch(context.TODO(), common.TaskKillerPrefix, clientv3.WithPrefix())
		for {
			select {
			case resp := <-watchChan:
				for _, v := range resp.Events {
					if v.Type == clientv3.EventTypePut {
						// 任务强杀不需要Value，只需Key(即taskName)
						name := strings.TrimPrefix(string(v.Kv.Key), common.TaskKillerPrefix)
						fmt.Printf("Kill %s\n", name)
						w.Scheduler.PushEvent(Event{Type: EventKill, Task: &Task{Name: name}})
					}
				}
			}
		}
	}()
}
