package worker

import (
	"CronJob/common"
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

type LogSink struct {
	client         *mongo.Client
	LogCollection  *mongo.Collection
	LogChan        chan *common.JobLog
	AutoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	_, err := logSink.LogCollection.InsertMany(context.TODO(), batch.Logs)
	fmt.Println("插入mongo错误:", err)
}

// 日志存储协程
func (logSink *LogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch // 当前的批次
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch
	)
	for {
		select {
		case log = <-logSink.LogChan:
			// 每次插入需要等待mongodb的一次请求往返，耗时可能因为网络慢花费比较长的时间
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// 让这个批次超时自动提交(给1秒时间)
				commitTimer = time.AfterFunc(time.Duration(Global.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.AutoCommitChan <- batch
						}
					}(logBatch))
			}
			// 把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, log)

			// 如果批次满了，就立即发送
			if len(logBatch.Logs) >= Global.JobLogBatchSize {
				logSink.saveLogs(logBatch)
				// 清空logbatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.AutoCommitChan:
			// 判断过期批次是否仍旧是当前的批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			logSink.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)
	if client, err = mongo.Connect(context.TODO(), &options.ClientOptions{Hosts: []string{Global.MongodbUri}}); err != nil {
		return
	}
	G_logSink = &LogSink{
		client:         client,
		LogCollection:  client.Database("cron").Collection("log"),
		LogChan:        make(chan *common.JobLog, 1000),
		AutoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go G_logSink.writeLoop()
	return
}

// 发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.LogChan <- jobLog:
	default:
		// 队列满了就丢弃
	}
}
