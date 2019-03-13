package thrift

import (
	"CronJob/common"
	"CronJob/master"
	"CronJob/master/rpc"
	"context"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type JobServer struct {
}

func (j *JobServer) ListJob(ctx context.Context) (r *rpc.Result_, err error) {
	var (
		jobList []*common.Job
		job     *common.Job
	)
	r = &rpc.Result_{}
	if jobList, err = master.G_jobMgr.ListJob(); err != nil {
		return
	}
	for _, job = range jobList {
		dataJob := &rpc.DataJob{
			Name:    job.Name,
			Command: job.Command,
			Expr:    job.CronExpr,
		}
		r.Job = append(r.Job, dataJob)
	}
	return
}

func InitRpcServer() {
	var (
		err       error
		transport *thrift.TServerSocket
	)
	if transport, err = thrift.NewTServerSocket(":9000"); err != nil {
		panic(err)
	}

	handler := &JobServer{}
	processor := rpc.NewJobServiceProcessor(handler)
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()
	server := thrift.NewTSimpleServer4(
		processor,
		transport,
		transportFactory,
		protocolFactory,
	)

	if err := server.Serve(); err != nil {
		panic(err)
	}
}
