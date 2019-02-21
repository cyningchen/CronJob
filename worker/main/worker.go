package main

import (
	"CronJob/worker"
	"fmt"
	"time"
)

func main() {
	var (
		err error
	)
	worker.InitArgs()
	if err = worker.InitConfig(worker.ConfFile); err != nil {
		fmt.Println("init config failed, ", err)
	}

	// 启动日志协程
	if err = worker.InitLogSink(); err != nil {
		fmt.Println("init logsink failed, ", err)
	}

	// 启动执行器
	if err = worker.InitExecutor(); err != nil {
		fmt.Println("init executor failed, ", err)
	}

	// 启动调度器
	if err = worker.InitScheduler(); err != nil {
		fmt.Println("init scheduler failed, ", err)
	}

	// 任务管理器
	if err = worker.InitJobMgr(); err != nil {
		fmt.Println("init job mgr failed, ", err)
	}

	for {
		time.Sleep(5 * time.Second)
	}
}
