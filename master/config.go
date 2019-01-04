package master

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
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	Webroot         string   `json:"webroot"`
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
	flag.StringVar(&ConfFile, "config", "./master.json", "指定配置文件")
	flag.Parse()
}
