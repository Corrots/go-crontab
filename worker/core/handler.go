package core

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func (s *Scheduler) eventHandler(e *Event) {
	taskName := e.Task.Name
	switch e.Type {
	case EventPut:
		expr, err := cronexpr.Parse(e.Task.Expression)
		if err != nil {
			fmt.Printf("parse cron expr err: %v\n", err)
			return
		}
		s.PlanTable[e.Task.Name] = &Plan{
			Task:     e.Task,
			Expr:     expr,
			NextTime: expr.Next(time.Now()),
		}
	case EventDelete:
		if _, ok := s.PlanTable[taskName]; ok {
			s.rwmutex.Lock()
			delete(s.PlanTable, taskName)
			s.rwmutex.Unlock()
		}
	case EventKill:
		// handle kill event
		if exec, ok := s.ExecTable[taskName]; ok {
			exec.CancelFunc()
		}
	}
}

func (s *Scheduler) resultHandler(res *Result) {
	l := &Log{
		TaskName:     res.TaskName,
		Command:      res.Command,
		Output:       string(res.Output),
		PlanTime:     getTimestamp(res.Exec.PlanTime),
		ScheduleTime: getTimestamp(res.Exec.ActualTime),
		StartTime:    getTimestamp(res.StartTime),
		EndTime:      getTimestamp(res.EndTime),
	}
	if res.Err != nil {
		l.Error = res.Err.Error()
	}
	s.LogProcessor.append(l)
}

func getTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1000 / 1000
}
