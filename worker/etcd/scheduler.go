package etcd

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/corrots/go-crontab/worker/utils"
)

var (
	timeZero = time.Time{}
)

type Scheduler struct {
	EventChan  chan Event
	PlanTable  map[string]*Plan
	ExecTable  map[string]*Exec
	ResultChan chan Result
	rwmutex    sync.RWMutex
}

func NewScheduler() Scheduler {
	return Scheduler{
		EventChan:  make(chan Event, 1000),
		PlanTable:  make(map[string]*Plan),
		ExecTable:  make(map[string]*Exec),
		ResultChan: make(chan Result, 1000),
		rwmutex:    sync.RWMutex{},
	}
}

func (s *Scheduler) PushEvent(e Event) {
	s.EventChan <- e
}

func (s *Scheduler) Run() {
	interval := s.getInterval()
	timer := time.NewTimer(interval)
	for {
		select {
		case e := <-s.EventChan:
			s.EventHandler(&e)
		case <-timer.C:
		case res := <-s.ResultChan:
			if res.Err != nil {
				log.Printf("exec task {%s} err: %v\n", res.TaskName, res.Err)
				continue
			}
			// 将res写入mongodb
			//if err := storeRes(&res); err != nil {
			//	fmt.Println(err)
			//	continue
			//}
			spent := res.EndTime.Sub(res.StartTime).Milliseconds()
			fmt.Printf("task {%v}, spent %v ms, output: %s\n", res.TaskName, spent, res.Output)
		}
		// 重新获取下次任务的间隔时间，并重置timer
		interval = s.getInterval()
		timer.Reset(interval)
	}
}

func storeRes(res *Result) error {
	mongo, err := utils.NewMongo()
	if err != nil {
		return fmt.Errorf("init mongo err: %v\n", err)
	}
	log := &utils.Log{
		TaskName: res.TaskName,
		//Command:  res.Output,
		//Error:  res.Err.Error(),
		Output: string(res.Output),
		//PlanTime:     "",
		//ScheduleTime: "0",
		StartTime: res.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:   res.EndTime.Format("2006-01-02 15:04:05"),
	}
	if res.Err != nil {
		log.Error = res.Err.Error()
	}
	return mongo.InsertLog(log)
}

func (s *Scheduler) getInterval() time.Duration {
	if len(s.PlanTable) == 0 {
		return time.Second
	}
	now := time.Now()
	var nearestTime time.Time
	for _, plan := range s.PlanTable {
		if plan.NextTime.Unix() <= now.Unix() {
			// @TODO task已到执行时间，尝试执行
			s.execute(plan)
			plan.NextTime = plan.Expr.Next(now)
		}
		// 获取接下来最接近的要执行的任务时间
		if nearestTime == timeZero || plan.NextTime.Before(nearestTime) {
			nearestTime = plan.NextTime
		}
	}
	return nearestTime.Sub(now)
}

func (s *Scheduler) execute(p *Plan) {
	taskName := p.Task.Name
	if _, existed := s.ExecTable[taskName]; existed {
		//fmt.Printf("previous task {%s} is running\n", taskName)
		return
	}

	taskExec := buildTaskExec(p)
	s.ExecTable[taskName] = taskExec
	go func(e *Exec, scheduler *Scheduler) {
		// 分布式锁
		lock, err := NewLock()
		if err != nil {
			fmt.Printf("new Lock err: %v\n", err)
			return
		}
		start := time.Now()
		if err := lock.Lock(taskExec.TaskName); err != nil {
			scheduler.ResultChan <- Result{
				TaskName:  taskExec.TaskName,
				Err:       err,
				StartTime: start,
				EndTime:   time.Now(),
			}
			return
		} else {
			start = time.Now()
			cmd := exec.CommandContext(taskExec.Ctx, "/bin/bash", "-c", p.Task.Command)
			output, err := cmd.CombinedOutput()
			scheduler.ResultChan <- Result{
				TaskName:  taskExec.TaskName,
				Output:    output,
				Err:       err,
				StartTime: start,
				EndTime:   time.Now(),
			}
		}
		defer lock.UnLock()
		// 将执行完成的Task移除
		scheduler.rwmutex.Lock()
		delete(scheduler.ExecTable, taskExec.TaskName)
		scheduler.rwmutex.Unlock()
	}(taskExec, s)
}

func buildTaskExec(plan *Plan) *Exec {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	return &Exec{
		TaskName:   plan.Task.Name,
		PlanTime:   plan.NextTime,
		ActualTime: time.Now(),
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}
}
