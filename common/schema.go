package common

const (
	TaskNamePrefix   = "/cron/jobs/"
	TaskKillerPrefix = "/cron/killer/"
	TaskLockPrefix   = "/cron/lock1/"
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
