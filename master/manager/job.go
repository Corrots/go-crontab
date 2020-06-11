package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/corrots/go-crontab/common/model"
)

const (
	JobNamePrefix   = "/cron/jobs/"
	JobKillerPrefix = "/cron/killer/"
)

func (jm *JobManager) SaveJob(job *model.Job) (prevJob *model.Job, err error) {
	key := getJobKey(job.Name)
	val, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("marshal job data err: %v\n", err)
	}
	putResp, err := jm.kv.Put(context.TODO(), key, string(val), clientv3.WithPrevKV())
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

func (jm *JobManager) GetJobByName(jobName string) (job *model.Job, err error) {
	key := getJobKey(jobName)
	getResp, err := jm.kv.Get(context.TODO(), key)
	if err != nil || getResp.Count == 0 {
		return nil, fmt.Errorf("etcd get err: %v\n", err)
	}
	if err := json.Unmarshal(getResp.Kvs[0].Value, &job); err != nil {
		return nil, fmt.Errorf("unmarshal job data err: %v\n", err)
	}
	return job, nil
}

func (jm *JobManager) GetJobs() (jobs []model.Job, err error) {
	getResp, err := jm.kv.Get(context.TODO(), JobNamePrefix, clientv3.WithPrefix())
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

func (jm *JobManager) DeleteJobs(jobName string) (prevJob *model.Job, err error) {
	key := getJobKey(jobName)
	delResp, err := jm.kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
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
func (jm *JobManager) JobKiller(jobName string) error {
	key := JobKillerPrefix + jobName
	// 为key设置lease
	ctx := context.Background()
	lease, err := jm.lease.Grant(ctx, 1)
	if err != nil {
		return fmt.Errorf("lease grant err: %v\n", err)
	}
	_, err = jm.kv.Put(ctx, key, "", clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("etcd put err: %v\n", err)
	}
	return nil
}

func getJobKey(jobName string) string {
	return JobNamePrefix + jobName
}
