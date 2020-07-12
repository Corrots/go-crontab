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
		TaskName:  res.TaskName,
		Command:   res.Command,
		Output:    string(res.Output),
		StartTime: string(res.StartTime.UnixNano() / 1000 / 1000),
		EndTime:   string(res.EndTime.UnixNano() / 1000 / 1000),
	}
	if res.Err != nil {
		l.Error = res.Err.Error()
	}
	s.LogProcessor.append(l)
}
