package core

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

type LogBash struct {
	Logs []interface{}
}

type LogSink struct {
	//AutoCommitChan chan LogBash
	LogsChan chan *Log
	Stack    []interface{}
	Limiter  int
}

func newLogSink() *LogSink {
	return &LogSink{
		LogsChan: make(chan *Log, 1000),
		Limiter:  100,
	}
}

func (ls *LogSink) Consume() {
	for {
		select {
		case log := <-ls.LogChan:

		}
	}
}
