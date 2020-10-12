package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var RespBadRequest = Response{Code: "-1", Message: "parse request data failed"}
var JobSaveFailed = Response{Code: "10002", Message: "cron job save failed"}
var JobLoadFailed = Response{Code: "10003", Message: "cron job load failed"}
var JobListEmpty = Response{Code: "10001", Message: "jobs list is empty"}
var JobDeleteFailed = Response{Code: "10005", Message: "cron job delete failed"}
var JobKillFailed = Response{Code: "10006", Message: "kill cron job failed"}

// log
var LogListFailed = Response{Code: "10007", Message: "get log list failed"}
var NoLog = Response{Code: "10008", Message: "log list is empty"}

//
var WorkerListFailed = Response{Code: "10009", Message: "get worker list failed"}

func sendResponse(c *gin.Context, statusCode int, resp *Response) {
	c.JSON(statusCode, gin.H{
		"code":    resp.Code,
		"message": resp.Message,
		"data":    resp.Data,
	})
	os.Exit(1)
}

func sendOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    "00",
		"message": "succeed",
		"data":    data,
	})
}
