package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/corrots/go-crontab/common/model"
)

const (
	PutEvent      = 1
	DeleteEvent   = 2
	JobNamePrefix = "/cron/jobs/"
)

type JobEvent struct {
	// event type 1: PUT; 2: DELETE
	Type int
	// event job info
	Job model.Job
}

func (jm *JobManager) createLock(jobName string) *Lock {
	return NewLock(jobName, jm.kv, jm.lease)
}

// 监听任务
func (jm *JobManager) WatchJobs() error {
	getResp, err := jm.kv.Get(context.TODO(), JobNamePrefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("etcd get err: %v\n", err)
	}
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("job list is empty")
	}
	for _, v := range getResp.Kvs {
		var job model.Job
		if err := json.Unmarshal(v.Value, &job); err != nil {
			return fmt.Errorf("job unmarshal err: %v\n", err)
		}
		// @TODO 将job推送至scheduler
		Scheduler.pushEvent(JobEvent{Type: PutEvent, Job: job})
	}
	// goroutine监听job更新
	go func() {
		latestRevision := getResp.Header.Revision + 1
		watchChan := jm.watcher.Watch(context.TODO(), JobNamePrefix, clientv3.WithRev(latestRevision), clientv3.WithPrefix())
		for {
			select {
			case resp := <-watchChan:
				for _, e := range resp.Events {
					var jobEvent JobEvent
					switch e.Type {
					case mvccpb.PUT:
						//fmt.Printf("PUT %s=%s\n", e.Kv.Key, e.Kv.Value)
						var job model.Job
						if err := json.Unmarshal(e.Kv.Value, &job); err != nil {
							log.Printf("unmarshal %s: %s err: %v\n", e.Kv.Key, e.Kv.Value, err)
							continue
						}
						// 构造更新event
						jobEvent = JobEvent{Type: PutEvent, Job: job}
					case mvccpb.DELETE:
						//fmt.Printf("DELETE Revision: %d\n", e.Kv.ModRevision)
						jobName := strings.TrimPrefix(string(e.Kv.Key), JobNamePrefix)
						// 构造删除event
						jobEvent = JobEvent{Type: DeleteEvent, Job: model.Job{Name: jobName}}
					}
					// @TODO 将event推送至scheduler
					Scheduler.pushEvent(jobEvent)
				}
			}
		}
	}()
	return nil
}
