package etcd

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/corrots/go-crontab/common"
)

type Lock struct {
	KV         clientv3.KV
	Lease      clientv3.Lease
	LeaseID    clientv3.LeaseID
	Ctx        context.Context
	CancelFunc context.CancelFunc
	Locked     bool
}

var ErrorLockOccupied = errors.New("锁已被占用")

func NewLock() (*Lock, error) {
	c, err := newEtcd()
	if err != nil {
		return nil, fmt.Errorf("etcd init err: %v\n", err)
	}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	return &Lock{
		KV:         clientv3.NewKV(c),
		Lease:      clientv3.NewLease(c),
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}, nil
}

// 1. 创建租约，并自动续租
// 2. 创建事务，事务抢锁
// 3. 成功返回，失败则释放租约
func (l *Lock) Lock(taskName string) error {
	// 创建租约
	lease, err := l.Lease.Grant(l.Ctx, 1)
	if err != nil {
		return fmt.Errorf("lease grant err: %v\n", err)
	}
	// 续租
	l.LeaseID = lease.ID
	aliveChan, err := l.Lease.KeepAlive(l.Ctx, l.LeaseID)
	if err != nil {
		return fmt.Errorf("lease keep alive err: %v\n", err)
	}
	// 监控自动续租resp chan
	go func() {
		for {
			if _, ok := <-aliveChan; !ok {
				return
			}
		}
	}()
	// 创建事务，抢锁
	txn := l.KV.Txn(context.TODO())
	key := common.TaskLockPrefix + taskName
	cmp := clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
	txn.If(cmp).Then(clientv3.OpPut(key, "", clientv3.WithLease(l.LeaseID))).Else(clientv3.OpGet(key))
	txnResp, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("txn commit err: %v\n", err)
	}
	if !txnResp.Succeeded {
		//kv := txnResp.Responses[0].GetResponseRange().Kvs[0]
		//return fmt.Errorf("{%s}抢锁失败, {%s}占用中\n", kv.Key, kv.Value)
		return ErrorLockOccupied
	}
	//fmt.Println("抢锁成功")
	return nil
}

func (l *Lock) UnLock() {
	l.Lease.Revoke(l.Ctx, l.LeaseID)
	l.CancelFunc()
}
