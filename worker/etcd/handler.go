package etcd

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func (s *Scheduler) EventHandler(e *Event) {
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
		//fmt.Printf("task: %+v\n", s.PlanTable[e.Task.Name].Task)
	case EventDelete:
		if _, ok := s.PlanTable[e.Task.Name]; ok {
			delete(s.PlanTable, e.Task.Name)
		}
	case EventKill:
		// @TODO handle kill event
	}
}
