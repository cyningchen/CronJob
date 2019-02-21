package worker

import (
	"CronJob/common"
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
)

// 注册节点到etcd /cron/workers/ip
type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIp string
}

var (
	G_Register *Register
)

func getLocalIp() (ipv4 string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logs.Error("get localIp failed, ", err)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4 = ipnet.IP.String()
				return
			}
		}
	}
	err = common.ERR_NO_LOCAL_IP_FOUNDED
	return
}

func (register *Register) KeepOnline() {
	var (
		err            error
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)
	cancelFunc = nil
	regKey = common.JOB_WORKER_DIR + register.localIp
	// 创建租约
	for {
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		// 自动续约
		leaseId = leaseGrantResp.ID
		if keepRespChan, err = register.lease.KeepAlive(context.TODO(), leaseId); err != nil {
			goto RETRY
		}
		cancelCtx, cancelFunc = context.WithCancel(context.TODO())
		// 注册到etcd
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseId)); err != nil {
			goto RETRY
		}
		// 处理续租应答协程
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto RETRY
				}
			}
		}
	RETRY:
		if cancelFunc != nil {
			cancelFunc()
		}
		time.Sleep(1 * time.Second)
	}

}

func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)
	config = clientv3.Config{
		Endpoints:   Global.EtcdEndpoints,
		DialTimeout: time.Duration(Global.EtcdDialTimeout) * time.Millisecond,
	}

	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("new etcd client failed: ", err)
		return
	}
	if localIp, err = getLocalIp(); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_Register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIp: localIp,
	}

	// 服务注册
	go G_Register.KeepOnline()
	return
}
