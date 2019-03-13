package main

import (
	"CronJob/master"
	"CronJob/master/thrift"
	"fmt"
	"time"
)

func main() {
	var (
		err error
	)
	master.InitArgs()
	if err = master.InitConfig(master.ConfFile); err != nil {
		fmt.Println("init config failed, ", err)
	}

	//  初始化服务发现
	if err = master.InitWorkerMgr(); err != nil {
		fmt.Println("init worker mgr failed, ", err)
	}

	// 日志管理器
	if err = master.InitLogMgr(); err != nil {
		fmt.Println("init log mgr failed, ", err)
	}

	// 任务管理器
	if err = master.InitJobMgr(); err != nil {
		fmt.Println("init job mgr failed, ", err)
	}

	// thrift rpc
	go thrift.InitRpcServer()

	if err := master.InitApiServer(); err != nil {
		fmt.Println("init apiserver failed: ", err)
	}
	for {
		time.Sleep(5 * time.Second)
	}

}
