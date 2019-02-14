package worker

import (
	"CronJob/common"
	"context"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd    *exec.Cmd
			err    error
			output []byte
			result *common.JobExecuteResult
		)

		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}
		result.StartTime = time.Now()
		// 执行shell命令
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		output, err = cmd.CombinedOutput()
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err

		// 任务执行完成后，把执行的结果返回给scheduler。 scheduler从executingTable中删除执行记录。
		G_scheduler.PushJobResult(result)
	}()
}

// 初始化调度器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}

// 回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.JobResultChan <- jobResult
}
