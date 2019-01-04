package master

import (
	"net/http"
	"net"
	"strconv"
	"fmt"
	"CronJob/common"
	"encoding/json"
	"time"
)

var (
	// 单例对象
	G_apiServer *ApiServer
)

// 任务的HTTP接口
type ApiServer struct {
	HttpServer *http.Server
}

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(Global.ApiPort)); err != nil {
		fmt.Printf("listen failed: %v", err)
		return
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(Global.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Global.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{
		HttpServer: httpServer,
	}
	go httpServer.Serve(listener)
	fmt.Println("http server is running....")
	return
}

// 保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		postjob string
		job     common.Job
		oldJob  *common.Job
		bytes []byte
	)
	if r.ParseForm(); err != nil{
		fmt.Printf("parse form failed: %v\n", err)
		return
	}
	postjob = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postjob), &job); err != nil{
		fmt.Printf("unmarshal postjob failed: %v\n", err)
		return
	}

	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil{
		fmt.Printf("save job failed: %v\n", err)
		return
	}

	if bytes, err = json.Marshal(oldJob); err != nil{
		fmt.Printf("marshal oldjob failed: %v\n", err)
		return
	}
	common.SendResponse(w,string(bytes),http.StatusCreated)
}
