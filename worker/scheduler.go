package worker

import (
	"CronJob/common"
	"fmt"
	"time"
)

// 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent               //  etcd任务事件队列
	jobPlanTable      map[string]*common.JobSchedulerPlan // 任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo   // 任务执行表
	JobResultChan     chan *common.JobExecuteResult       // 任务执行结果队列
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulerPlan
		jobExisted      bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: // 保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulerPlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: // 删除任务事件
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// 处理任务结果
func (scheduler *Scheduler) handleJobResult(jobResult *common.JobExecuteResult) {
	// 删除执行状态
	delete(scheduler.jobExecutingTable, jobResult.ExecuteInfo.Job.Name)

	fmt.Println("任务执行完成", jobResult.ExecuteInfo.Job.Name, string(jobResult.Output), jobResult.Err)
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	// 调度和执行是两件事
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)

	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println("尚未退出，跳过执行:", jobPlan.Job.Name)
		return
	}
	// 构建执行状态信息
	jobExecuteInfo = common.BuildJonExecuteInfo(jobPlan)

	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	fmt.Println("执行任务", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	G_executor.ExecuteJob(jobExecuteInfo)
	return
}

// 重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulerPlan
		now      time.Time
		nearTime *time.Time
	)

	// 如果任务表为空话，随便睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	// 当前时间
	now = time.Now()

	// 遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}

		// 统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔（最近要执行的任务调度时间 - 当前时间）
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// 调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)

	// 初始化一次(1秒)
	scheduleAfter = scheduler.TrySchedule()

	// 调度的延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	// 定时任务common.Job
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //监听任务变化事件
			// 对内存中维护的任务列表做增删改查
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: // 最近的任务到期了
			// 调度一次任务
			scheduleAfter = scheduler.TrySchedule()
			// 重置调度间隔
			scheduleTimer.Reset(scheduleAfter)
		case jobResult = <-scheduler.JobResultChan:
			scheduler.handleJobResult(jobResult)
		}
	}
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// 初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		JobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}
