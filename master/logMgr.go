package master

import (
	"CronJob/common"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type LogMgr struct {
	client        *mongo.Client
	LogCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client *mongo.Client
	)
	if client, err = mongo.Connect(context.TODO(), &options.ClientOptions{Hosts: []string{Global.MongodbUri}}); err != nil {
		return
	}
	G_logMgr = &LogMgr{
		client:        client,
		LogCollection: client.Database("cron").Collection("log"),
	}
	return
}

func (logMgr *LogMgr) ListLog(name string, skip int64, limit int64) (logArr []*common.JobLog, err error) {
	var (
		filter  *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor  *mongo.Cursor
		jobLog  *common.JobLog
	)

	logArr = make([]*common.JobLog, 0)

	filter = &common.JobLogFilter{JobName: name}
	logSort = &common.SortLogByStartTime{SortOrder: -1}
	if cursor, err = logMgr.LogCollection.Find(context.TODO(), filter, &options.FindOptions{Sort: logSort, Skip: &skip, Limit: &limit}); err != nil {
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		if err = cursor.Decode(jobLog); err != nil {
			continue // 日志不合法
		}
		logArr = append(logArr, jobLog)
	}
	return
}
