package manager

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
)

const (
	LockNamePrefix = "/cron/lock/"
)

var ErrorLockFailed = errors.New("抢锁失败")

type Lock struct {
	JobName    string
	KV         clientv3.KV
	Lease      clientv3.Lease
	LeaseID    clientv3.LeaseID
	CancelFunc context.CancelFunc
	IsLocked   bool
}

func NewLock(jobName string, kv clientv3.KV, lease clientv3.Lease) *Lock {
	return &Lock{
		JobName: jobName,
		KV:      kv,
		Lease:   lease,
	}
}

// 1.创建租约，并自动续租
// 2.通过事务抢锁
// 3.抢锁成功返回，失败则释放租约
func (l *Lock) TryLock() error {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	lease, err := l.Lease.Grant(ctx, 1)
	if err != nil {
		return fmt.Errorf("lease grant err: %v\n", err)
	}
	aliveChan, err := l.Lease.KeepAlive(ctx, lease.ID)
	if err != nil {
		return fmt.Errorf("keep lease alive err: %v\n", err)
	}
	// 消费自动续租resp
	go func() {
		for {
			if _, ok := <-aliveChan; !ok {
				return
			}
		}
	}()

	txn := l.KV.Txn(context.TODO())
	key := LockNamePrefix + l.JobName
	cmp := clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
	txn.If(cmp).Then(clientv3.OpPut(key, "", clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(key))
	txnResp, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("txn commit err: %v\n", err)
	}
	if txnResp.Succeeded {
		l.CancelFunc = cancelFunc
		l.LeaseID = lease.ID
		l.IsLocked = true
		return nil
	}
	log.Printf("抢锁失败, %s\n", txnResp.Responses[0].GetResponseRange().Kvs[0].Key)
	return ErrorLockFailed
}

func (l *Lock) Unlock() {
	if l.IsLocked {
		l.CancelFunc()
		l.Lease.Revoke(context.TODO(), l.LeaseID)
	}
}
