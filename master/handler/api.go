package handler

import (
	"log"
	"net/http"

	"github.com/corrots/go-crontab/common/model"
	"github.com/corrots/go-crontab/master/manager"
	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) {
	var job model.Job
	if err := c.Bind(&job); err != nil {
		log.Printf("bind request body err: %v\n", err)
		sendResponse(c, http.StatusBadRequest, &RespBadRequest)
		return
	}
	if err := job.Validation(); err != nil {
		sendResponse(c, http.StatusBadRequest, &RespBadRequest)
		return
	}
	prevJob, err := manager.JM.SaveJob(&job)
	if err != nil {
		log.Println(err)
		sendResponse(c, http.StatusInternalServerError, &JobSaveFailed)
		return
	}
	sendOK(c, prevJob)
}

func List(c *gin.Context) {
	jobs, err := manager.JM.GetJobs()
	if err != nil {
		log.Println(err)
		sendResponse(c, http.StatusInternalServerError, &JobListFailed)
		return
	}
	sendOK(c, jobs)
}

func Load(c *gin.Context) {
	jobName := c.Param("jobName")
	if jobName == "" {
		sendResponse(c, http.StatusBadRequest, &RespBadRequest)
		return
	}
	job, err := manager.JM.GetJobByName(jobName)
	if err != nil {
		log.Println(err)
		sendResponse(c, http.StatusInternalServerError, &JobLoadFailed)
		return
	}
	sendOK(c, job)
}

func Delete(c *gin.Context) {
	jobName := c.Param("jobName")
	if jobName == "" {
		sendResponse(c, http.StatusBadRequest, &RespBadRequest)
		return
	}
	prevJob, err := manager.JM.DeleteJob(jobName)
	if err != nil {
		log.Println(err)
		sendResponse(c, http.StatusInternalServerError, &JobDeleteFailed)
		return
	}
	sendOK(c, prevJob)
}

func Kill(c *gin.Context) {
	jobName := c.Param("jobName")
	if jobName == "" {
		sendResponse(c, http.StatusBadRequest, &RespBadRequest)
		return
	}
	if err := manager.JM.JobKiller(jobName); err != nil {
		sendResponse(c, http.StatusInternalServerError, &JobKillFailed)
		return
	}
	sendOK(c, nil)
}
