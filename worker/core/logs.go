package core

import (
	"log"
	"time"
)

type Log struct {
	// 任务名称
	TaskName string `bson:"task_name"`
	// shell指令
	Command string `bson:"command"`
	// shell执行error
	Error string `bson:"error"`
	// shell执行stdout
	Output string `bson:"output"`
	// 计划开始时间
	PlanTime string `bson:"plan_time"`
	// 实际调度时间
	ScheduleTime string `bson:"schedule_time"`
	// 任务执行开始时间
	StartTime string `bson:"start_time"`
	// 任务执行结束时间
	EndTime string `bson:"end_time"`
}

type LogBatch struct {
	Logs []interface{}
}

type LogProcessor struct {
	mongo      *Mongo
	logChan    chan *Log
	logBatches chan *LogBatch
	limiter    int
}

func NewLogProcessor() *LogProcessor {
	mongo, err := NewMongo()
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
					err := lp.mongo.InsertMany(logs)
					if err != nil {
						log.Printf("mongo insert logs err: %v\n", err)
					}
				}(batch.Logs)
				// 写入完成后清空bash中的logs
				batch.Logs = nil
				timer.Stop()
			}
		case timeoutBatch := <-lp.logBatches:
			if timeoutBatch != batch {
				continue
			}
			go func(logs []interface{}) {
				err := lp.mongo.InsertMany(logs)
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
