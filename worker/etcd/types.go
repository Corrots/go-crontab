package etcd

import (
	"context"
	"time"

	"github.com/gorhill/cronexpr"
)

type Task struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	Expression string `json:"expression"`
}

type Plan struct {
	Task     *Task
	Expr     *cronexpr.Expression
	NextTime time.Time
}

type Exec struct {
	TaskName   string
	PlanTime   time.Time
	ActualTime time.Time
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

type Result struct {
	TaskName  string
	Command   string
	Output    []byte
	Err       error
	StartTime time.Time
	EndTime   time.Time
}
