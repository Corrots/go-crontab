package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/corrots/go-crontab/common"
	"github.com/corrots/go-crontab/common/model"
)

func (m *Manager) SaveJob(job *model.Job) (prevJob *model.Job, err error) {
	key := getJobKey(job.Name)
	val, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("marshal job data err: %v\n", err)
	}
	putResp, err := m.kv.Put(context.TODO(), key, string(val), clientv3.WithPrevKV())
	if err != nil || putResp.Header.Revision == 0 {
		return nil, fmt.Errorf("etcd put err: %v\n", err)
	}
	// 如果是更新，返回PrevKv
	if putResp.PrevKv != nil {
		json.Unmarshal(putResp.PrevKv.Value, &prevJob)
		return
	}
	return
}

func (m *Manager) GetJobByName(jobName string) (job *model.Job, err error) {
	key := getJobKey(jobName)
	getResp, err := m.kv.Get(context.TODO(), key)
	if err != nil || getResp.Count == 0 {
		return nil, fmt.Errorf("etcd get err: %v\n", err)
	}
	if err := json.Unmarshal(getResp.Kvs[0].Value, &job); err != nil {
		return nil, fmt.Errorf("unmarshal job data err: %v\n", err)
	}
	return job, nil
}

func (m *Manager) GetJobs() (jobs []model.Job, err error) {
	getResp, err := m.kv.Get(context.TODO(), common.TaskNamePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd get err: %v\n", err)
	}
	if getResp.Count == 0 {
		return nil, errors.New("job list is empty")
	}
	for _, v := range getResp.Kvs {
		var job model.Job
		if err := json.Unmarshal(v.Value, &job); err != nil {
			return nil, fmt.Errorf("unmarshal job data err: %v\n", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (m *Manager) DeleteJob(jobName string) (prevJob *model.Job, err error) {
	key := getJobKey(jobName)
	delResp, err := m.kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
	if err != nil {
		return nil, fmt.Errorf("etcd delete err: %v\n", err)
	}
	if len(delResp.PrevKvs) == 0 {
		return nil, fmt.Errorf("ectd key: {%s} doesn't exist\n", jobName)
	}
	json.Unmarshal(delResp.PrevKvs[0].Value, prevJob)
	return prevJob, nil
}

// 更新 key = /cron/killer/jobName
func (m *Manager) JobKiller(jobName string) error {
	key := common.TaskKillerPrefix + jobName
	// 为key设置lease
	ctx := context.Background()
	lease, err := m.lease.Grant(ctx, 1)
	if err != nil {
		return fmt.Errorf("lease grant err: %v\n", err)
	}
	_, err = m.kv.Put(ctx, key, "", clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("etcd put err: %v\n", err)
	}
	return nil
}

func getJobKey(jobName string) string {
	return common.TaskNamePrefix + jobName
}
