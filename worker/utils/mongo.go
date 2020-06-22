package utils

type JobLog struct {
	JobName string `json:"job_name"`
	Command string `json:"command"`
	Error   string `json:"error"`
	Output  string `json:"output"`
	// 计划开始时间
	PlanTime int64 `json:"plan_time"`
	// 实际调度时间
	ScheduleTime int64 `json:"schedule_time"`
	// 任务执行开始时间
	StartTime int64 `json:"start_time"`
	// 任务执行结束时间
	EndTime int64 `json:"end_time"`
}
