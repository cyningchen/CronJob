package master

import (
	"CronJob/common"
	"CronJob/master/defs"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		httpServer    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)

	// 静态资源目录
	staticDir = http.Dir(Global.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	httpServer = &http.Server{
		Addr:         ":" + strconv.Itoa(Global.ApiPort),
		ReadTimeout:  time.Duration(Global.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Global.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	go httpServer.ListenAndServe()
	fmt.Println("http server is running....")
	return
}

// 保存任务接口
// POST /job/save job={json}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		postjob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	if r.ParseForm(); err != nil {
		common.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	postjob = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postjob), &job); err != nil {
		common.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	if bytes, err = json.Marshal(oldJob); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	common.SendResponse(w, string(bytes), http.StatusCreated)
}

// 删除任务接口
// POST /job/delete name=job
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldjob *common.Job
		bytes  []byte
	)
	if err = r.ParseForm(); err != nil {
		common.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	name = r.PostForm.Get("name")
	if oldjob, err = G_jobMgr.DelJob(name); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	if bytes, err = json.Marshal(oldjob); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	common.SendResponse(w, string(bytes), http.StatusOK)
}

// 列出所有任务接口
//
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobList []*common.Job
		bytes   []byte
	)
	if jobList, err = G_jobMgr.ListJob(); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
	}
	if bytes, err = json.Marshal(jobList); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
	}
	common.SendResponse(w, string(bytes), http.StatusOK)
}

func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		name string
	)
	if err = r.ParseForm(); err != nil {
		common.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	name = r.PostForm.Get("name")
	if err = G_jobMgr.KillJob(name); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	common.SendResponse(w, "success", http.StatusOK)
}

// 查询任务日志
func handleJobLog(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		name       string
		skipParam  string
		limitParam string
		skip       int64
		limit      int64
		logArr     []*common.JobLog
		bytes      []byte
	)
	if err = r.ParseForm(); err != nil {
		common.SendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	name = r.Form.Get("name")
	skipParam = r.Form.Get("skip")
	limitParam = r.Form.Get("limit")
	if skip, err = strconv.ParseInt(skipParam, 10, 64); err != nil {
		skip = 0
	}
	if limit, err = strconv.ParseInt(limitParam, 10, 64); err != nil {
		limit = 20
	}
	if logArr, err = G_logMgr.ListLog(name, skip, limit); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
		return
	}
	if bytes, err = json.Marshal(logArr); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
	}
	common.SendResponse(w, string(bytes), http.StatusOK)
}

func handleWorkerList(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		workerArr []string
		bytes     []byte
	)
	if workerArr, err = G_workerMgr.ListWorkers(); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
	}
	if bytes, err = json.Marshal(workerArr); err != nil {
		common.SendErrorResponse(w, defs.ErrorInternalFault)
	}
	common.SendResponse(w, string(bytes), http.StatusOK)
}
