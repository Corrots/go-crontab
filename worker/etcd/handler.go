package etcd

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func (s *Scheduler) EventHandler(e *Event) {
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
		// @TODO handle kill event
		if exec, ok := s.ExecTable[taskName]; ok {
			exec.CancelFunc()
		}
	}
}
