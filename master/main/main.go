package main

import (
	"CronJob/master"
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
	// 任务管理器
	if err = master.InitJobMgr(); err != nil {
		fmt.Println("init job mgr failed, ", err)
	}

	if err := master.InitApiServer(); err != nil {
		fmt.Println("init apiserver failed: ", err)
	}
	for{
		time.Sleep(5*time.Second)
	}
}
