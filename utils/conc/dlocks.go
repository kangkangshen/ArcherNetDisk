package conc

import (
	"context"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/etcd-io/etcd/lease"
	"github.com/kangkangshen/ArcherNetDisk/config"
	"github.com/rs/xid"
	"log"
	"math/rand"
	"sync"
	"time"
)

//分布式锁,非可重入,通道实现
type DLocker interface {
	sync.Locker                                                      //阻塞方式的锁接口。如果etcd不可用，抛出panic
	TryLockOnce() bool                                               //非阻塞方式的所接口
	TryLockWithTimeout(ctx context.Context, time time.Duration) bool //带超时时间的锁接口,可以取消
	MaybeUnlockAfter() time.Duration                                 //预估？？时间后所释放，当锁是租约机制
}

//etcd的默认实现版本
type etcdDLocker struct {
	internalLock sync.Mutex //need lock guarantee 'state' new for local other goroutines,and coordinate
	state        int
	client       *clientv3.Client //etcdClient
	target       string           //目标字段 需要显式传入
	//unlockSignalChan chan interface{}	//广播释放锁的信号
	lease          clientv3.Lease
	currentLeaseId clientv3.LeaseID
	ctx            context.Context
	cancelfunc     context.CancelFunc
	id             string
}

func NewDlocker(client *clientv3.Client, target string) DLocker {
	dLocker := &etcdDLocker{client: client, target: config.DLOCKER_PREFIX + target, lease: clientv3.NewLease(client)}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	dLocker.ctx = ctx
	dLocker.cancelfunc = cancelFunc
	dLocker.state = config.DLOCKER_NEW
	dLocker.id = config.DLOCKER_PREFIX + xid.New().String()
	//dLocker.unlockSignalChan<-config.SIGNAL_UNLOCK
	return dLocker
}

func (locker *etcdDLocker) Lock() {
	var (
		logger *log.Logger
		kv     clientv3.KV
		tx     clientv3.Txn
		lease  clientv3.Lease
		ctx    context.Context
		lResp  *clientv3.LeaseGrantResponse
		txResp *clientv3.TxnResponse
		err    error
		count  int
	)

	locker.internalLock.Lock()
	defer locker.internalLock.Unlock()
	kv = clientv3.NewKV(locker.client)
	logger = config.Logger
	lease = locker.lease
	ctx = locker.ctx
	if lResp, err = lease.Grant(ctx, config.DLOCK_LEASE_TIME); err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, err = lease.KeepAliveOnce(ctx, lResp.ID); err != nil {
					panic(err)
				}
				time.Sleep(time.Duration(config.DLOCK_LEASE_TIME))
			}
		}
	}()

	tx = kv.Txn(ctx)
	for txResp, err = tx.If(clientv3.Compare(clientv3.CreateRevision(locker.target), "=", 0)).
		Then(clientv3.OpPut(locker.target, locker.id, clientv3.WithLease(lResp.ID))).
		Commit(); err != nil || !txResp.Succeeded; {
		if err != nil {
			panic(err)
		}
		//抢占失败,停止随机时间，如果有的话,否则不停获取
		count++
		time.Sleep(time.Duration(rand.Int63n(int64(time.Second))))
		logger.Printf("failed to apply for lock, try %d times again\n", count)

	}
	//抢占成功
	if txResp.Succeeded {
		locker.state = config.DLOCER_LOCK
		logger.Printf("experienced %d times application locks successfully", count)
		return
	}
}

func (locker *etcdDLocker) Unlock() {
	var (
		logger    *log.Logger
		kv        clientv3.KV
		tx        clientv3.Txn
		txResp    *clientv3.TxnResponse
		err       error
		ctx       context.Context
		cacelFunc context.CancelFunc
	)

	locker.internalLock.Lock()
	defer locker.internalLock.Unlock()
	logger = config.Logger
	kv = clientv3.NewKV(locker.client)
	ctx = locker.ctx
	cacelFunc = locker.cancelfunc
	tx = kv.Txn(ctx)
	if txResp, err = tx.If(clientv3.Compare(clientv3.Value(locker.target), "=", locker.id)).
		Then(clientv3.OpDelete(locker.target)).
		Commit(); err != nil || !txResp.Succeeded {
		if err != nil {
			panic(err)
		}
		if !txResp.Succeeded {
			logger.Println("wrong lock state,expected lock to release successfully but not")
			panic("wrong lock state,expected lock to release successfully but not，Probably because the lock holder is not himself or the lock has not been held")
		}
	}
	cacelFunc()
	locker.state = config.DLOCKER_UNLOCK
}

func (locker *etcdDLocker) TryLockOnce() bool {
	var (
		logger *log.Logger
		kv     clientv3.KV
		tx     clientv3.Txn
		txResp *clientv3.TxnResponse
		err    error
	)

	locker.internalLock.Lock()
	defer locker.internalLock.Unlock()
	logger = config.Logger
	kv = clientv3.NewKV(locker.client)
	tx = kv.Txn(context.TODO())
	if txResp, err = tx.If(clientv3.Compare(clientv3.CreateRevision(locker.target), "=", 0)).
		Then(clientv3.OpPut(locker.target, "locked")).
		Commit(); err != nil || !txResp.Succeeded {
		if err != nil {
			panic(err)
		}
		if !txResp.Succeeded {
			logger.Println("failed to try to apply for lock")
			return false
		}
	}
	logger.Println("succeed to try to apply for lock")
	locker.state = config.DLOCER_LOCK
	return true
}

func (locker *etcdDLocker) TryLockWithTimeout(ctx context.Context, d time.Duration) bool {
	var (
		logger   *log.Logger
		kv       clientv3.KV
		tx       clientv3.Txn
		txResp   *clientv3.TxnResponse
		err      error
		count    int
		endT     time.Time
		currentT time.Time
	)
	endT = time.Now().Add(d)
	locker.internalLock.Lock()
	defer locker.internalLock.Unlock()
	logger = config.Logger
	kv = clientv3.NewKV(locker.client)
	tx = kv.Txn(ctx)
	for txResp, err = tx.If(clientv3.Compare(clientv3.CreateRevision(locker.target), "=", 0)).
		Then(clientv3.OpPut(locker.target, "locked")).
		Commit(); err != nil || !txResp.Succeeded; {
		if err != nil {
			panic(err)
		}
		if txResp.Succeeded {
			logger.Printf("experienced %d times application locks successfully", count)
			return true
		}
		select {
		case <-ctx.Done():
			return false //获取失败,被取消，不继续进行
		default:
			//抢占失败,停止随机时间，如果有的话,否则不停获取
			count++
			time.Sleep(time.Duration(rand.Int63n(int64(time.Second))))
			logger.Printf("failed to apply for lock, try %d times again\n", count)
			if currentT = time.Now(); currentT.After(endT) {
				return false
			}
		}
	}
	//抢占成功
	if txResp.Succeeded {
		logger.Printf("experienced %d times application locks successfully", count)
		locker.state = config.DLOCER_LOCK
		return true
	}
	return false

}

func (locker *etcdDLocker) LockerHolder() string {
	var (
		kv      clientv3.KV
		ctx     context.Context
		target  string
		err     error
		getResp *clientv3.GetResponse
	)
	kv = clientv3.NewKV(locker.client)
	ctx = locker.ctx
	target = locker.target
	if getResp, err = kv.Get(ctx, target); err != nil {
		panic(err)
	}
	if len(getResp.Kvs) == 0 {
		return ""
	} else {
		return string(getResp.Kvs[0].Value)
	}
}

func (locker *etcdDLocker) MaybeUnlockAfter() time.Duration {
	var (
		err error
		ltr *clientv3.LeaseTimeToLiveResponse
	)
	if ltr, err = locker.lease.TimeToLive(context.TODO(), locker.currentLeaseId); err != nil {
		if err == lease.ErrLeaseNotFound {
			return 0
		} else {
			panic(err)
		}
	}
	return time.Duration(ltr.TTL)

}

/*
func (locker *etcdDLocker) Release() error{
	var (
		kv clientv3.KV
		err error
		tx clientv3.Txn
		txResp *clientv3.TxnResponse
		ctx context.Context
	)

	locker.internalLock.Lock()
	defer locker.internalLock.Unlock()
	//check state
	if(locker.state==config.DLOCKER_RELEASE){
		panic(errors.New("wrong locker state,locker has been released"))
	}
	ctx=locker.ctx
	kv=clientv3.NewKV(locker.client)
	tx = kv.Txn(ctx)
	for txResp, err = tx.If(clientv3.Compare(clientv3.CreateRevision(locker.target), "=", 0)).
		Then(clientv3.OpPut(locker.target, "locked")).
		Commit(); err != nil || !txResp.Succeeded; {

	}
	return nil;
}

*/
