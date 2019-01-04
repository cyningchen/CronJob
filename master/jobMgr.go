package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
	"encoding/json"
	"context"
	"CronJob/common"
)

var (
	G_jobMgr *JobMgr
)

// 任务管理
type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

func InitJobMgr() (err error)  {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease

	)
	config = clientv3.Config{
		Endpoints:Global.EtcdEndpoints,
		DialTimeout:time.Duration(Global.EtcdDialTimeout) * time.Millisecond,
	}

	//建立连接
	if client, err = clientv3.New(config); err != nil{
		fmt.Println("new etcd client failed: ", err)
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client:client,
		kv:kv,
		lease:lease,
	}
	return
}

func (j *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error)  {
	// 任务保存到/cron/jobs/任务名 ->json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	jobKey = "/cron/jobs/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil{
		return
	}
	// save to etcd
	if putResp, err = j.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV()); err != nil{
		return
	}
	if putResp.PrevKv != nil{
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil{
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}
