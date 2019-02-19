package worker

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var (
	Global   Config
	ConfFile string
)

// 程序配置
type Config struct {
	EtcdEndpoints       []string `json:"etcdEndpoints"`
	EtcdDialTimeout     int      `json:"etcdDialTimeout"`
	MongodbUri          string   `json:"mongodbUri"`
	MongodbTimeout      int      `json:"mongodbTimeout"`
	JobLogBatchSize     int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout int      `json:"jobLogCommitTimeout"`
}

func InitConfig(filename string) (err error) {
	var (
		content []byte
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = json.Unmarshal(content, &Global); err != nil {
		return
	}
	return
}

//解析命令行参数
func InitArgs() {
	flag.StringVar(&ConfFile, "config", "./worker.json", "指定配置文件")
	flag.Parse()
}
