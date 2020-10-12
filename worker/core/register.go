package core

import (
	"context"
	"fmt"

	"github.com/corrots/go-crontab/common"

	"github.com/coreos/etcd/clientv3"
)

const (
	DefaultLeaseTTL = 10
)

// 将当前worker注册到etcd中，并自动续租
type Register struct {
	Client  *clientv3.Client
	KV      clientv3.KV
	Lease   clientv3.Lease
	Watcher clientv3.Watcher
}

func NewRegister() (*Register, error) {
	c, err := newEtcd()
	if err != nil {
		return nil, fmt.Errorf("etcd init err: %v\n", err)
	}
	return &Register{
		Client:  c,
		KV:      clientv3.NewKV(c),
		Lease:   clientv3.NewLease(c),
		Watcher: clientv3.NewWatcher(c),
	}, nil
}

func (r *Register) KeepOnline() error {
	// 创建租约
	ctx, cancel := context.WithCancel(context.TODO())
	lease, err := r.Lease.Grant(ctx, DefaultLeaseTTL)
	if err != nil {
		return fmt.Errorf("lease grant err: %v\n", err)
	}
	// 自动续租
	aliveChan, err := r.Lease.KeepAlive(ctx, lease.ID)
	if err != nil {
		cancel()
		return fmt.Errorf("lease keep alive err: %v\n", err)
	}
	// 将当前worker注册到etcd
	workerID := common.WorkerListPrefix + getWorkerIP()
	if _, err := r.KV.Put(ctx, workerID, "", clientv3.WithLease(lease.ID)); err != nil {
		cancel()
		return fmt.Errorf("etcd put err: %v\n", err)
	}
	// 消费 <-chan *LeaseKeepAliveResponse
	go func() {
		for {
			if _, ok := <-aliveChan; !ok {
				return
			}
		}
	}()
	return nil
}
