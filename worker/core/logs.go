package core

import (
	"log"
	"time"

	"github.com/corrots/go-crontab/common"
)

type Log struct {
	// 任务名称
	TaskName string `json:"task_name";bson:"task_name"`
	// shell指令
	Command string `json:"command";bson:"command"`
	// shell执行error
	Error string `json:"error";bson:"error"`
	// shell执行stdout
	Output string `json:"output";bson:"output"`
	// 计划开始时间
	PlanTime int64 `json:"plan_time";bson:"plan_time"`
	// 实际调度时间
	ScheduleTime int64 `json:"schedule_time";bson:"schedule_time"`
	// 任务执行开始时间
	StartTime int64 `json:"start_time";bson:"start_time"`
	// 任务执行结束时间
	EndTime int64 `json:"end_time";bson:"end_time"`
}

type LogBatch struct {
	Logs []interface{}
}

type LogProcessor struct {
	mongo      *common.Mongo
	logChan    chan *Log
	logBatches chan *LogBatch
	limiter    int
}

func NewLogProcessor() *LogProcessor {
	mongo, err := common.NewMongo()
	if err != nil {
		log.Fatal(err)
	}
	return &LogProcessor{
		mongo:      mongo,
		logChan:    make(chan *Log, 1000),
		logBatches: make(chan *LogBatch, 1000),
		limiter:    2,
	}
}

func (lp *LogProcessor) Consumer() {
	var batch *LogBatch
	var timer *time.Timer
	for {
		select {
		case l := <-lp.logChan:
			if batch == nil {
				batch = &LogBatch{}
				timer = time.AfterFunc(time.Second*5, func(batch *LogBatch) func() {
					return func() {
						lp.logBatches <- batch
					}
				}(batch))
			}
			batch.Logs = append(batch.Logs, l)
			if len(batch.Logs) >= lp.limiter {
				go func(logs []interface{}) {
					err := lp.mongo.InsertLogs(logs)
					if err != nil {
						log.Printf("mongo insert logs err: %v\n", err)
					}
				}(batch.Logs)
				// 写入完成后清空bash中的logs
				batch.Logs = nil
				timer.Stop()
			}
		//	超过5s，即使batch中日志数量未超过10条，也会被提交保存至mongo
		case timeoutBatch := <-lp.logBatches:
			if timeoutBatch != batch {
				continue
			}
			go func(logs []interface{}) {
				err := lp.mongo.InsertLogs(logs)
				if err != nil {
					log.Printf("mongo insert logs err: %v\n", err)
				}
			}(timeoutBatch.Logs)
			batch.Logs = nil
		}
	}
}

func (lp *LogProcessor) append(l *Log) {
	select {
	case lp.logChan <- l:
	default:
	}
}
