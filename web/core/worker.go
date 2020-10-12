package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/corrots/go-crontab/common"
)

func (m *Manager) GetWorkerList() ([]string, error) {
	ctx := context.TODO()
	getResp, err := m.kv.Get(ctx, common.WorkerListPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd get err: %kv\n", err)
	}
	var workers []string
	for _, kv := range getResp.Kvs {
		workers = append(workers, extractWorkerIP(kv.Key))
	}
	return workers, nil
}

func extractWorkerIP(key []byte) string {
	return strings.TrimLeft(string(key), common.WorkerListPrefix)
}
