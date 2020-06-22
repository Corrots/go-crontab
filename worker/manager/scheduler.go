package manager

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/corrots/go-crontab/worker/utils"

	"github.com/corrots/go-crontab/common/model"
	"github.com/gorhill/cronexpr"
)

type JobScheduler struct {
	EventChan  chan JobEvent
	Plans      map[string]*JobPlan
	Executions map[string]*JobExecution
	ResChan    chan ExecRes
}

type ExecRes struct {
	Execution *JobExecution
	Output    []byte
	Error     error
	StartTime time.Time
	EndTime   time.Time
}

type JobExecution struct {
	Name       string
	PlanTime   time.Time
	ExecTime   time.Time
	CancelCtx  context.Context
	CancelFunc context.CancelFunc
}

type JobPlan struct {
	Job        *model.Job
	Expression *cronexpr.Expression
	NextTime   time.Time
}

var Scheduler *JobScheduler

func InitScheduler() {
	Scheduler = &JobScheduler{
		EventChan:  make(chan JobEvent, 1000),
		Plans:      make(map[string]*JobPlan),
		Executions: make(map[string]*JobExecution),
		ResChan:    make(chan ExecRes, 1000),
	}
	go Scheduler.run()
}

func (s *JobScheduler) execute(plan *JobPlan) {
	name := plan.Job.Name
	if _, exist := s.Executions[name]; exist {
		fmt.Printf("previous task {%s} is running\n", name)
		return
	}

	je := buildJobExecution(plan)
	s.Executions[name] = je

	go func(e *JobExecution, job *model.Job) {
		var res ExecRes
		// 随机sleep, 打散任务
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		// 获取分布式锁
		startTime := time.Now()
		lock := JM.createLock(job.Name)
		if err := lock.TryLock(); err != nil {
			res = ExecRes{
				Execution: e,
				Error:     err,
				StartTime: startTime,
				EndTime:   time.Now(),
			}
		} else {
			startTime = time.Now()
			cmd := exec.CommandContext(je.CancelCtx, "/bin/bash", "-c", job.Command)
			output, err := cmd.CombinedOutput()
			//fmt.Printf("cmd exec err: %v\n", err)
			res = ExecRes{
				Execution: e,
				Output:    output,
				Error:     err,
				StartTime: startTime,
				EndTime:   time.Now(),
			}
		}
		s.ResChan <- res
		defer lock.Unlock()
	}(je, plan.Job)

	if _, exist := s.Executions[name]; exist {
		delete(s.Executions, name)
	}
}

func buildJobExecution(plan *JobPlan) *JobExecution {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	return &JobExecution{
		Name:       plan.Job.Name,
		PlanTime:   plan.NextTime,
		ExecTime:   time.Now(),
		CancelCtx:  ctx,
		CancelFunc: cancelFunc,
	}
}

func (s *JobScheduler) run() {
	after := s.tryScheduler()
	timer := time.NewTimer(after)
	for {
		select {
		case event := <-s.EventChan:
			s.eventHandler(event)
		case <-timer.C:
		case res := <-s.ResChan:
			if res.Error != nil {
				if res.Error != ErrorLockFailed {
					jobLog := &utils.JobLog{
						JobName:      res.Execution.Name,
						Command:      nil,
						Error:        "",
						Output:       string(res.Output),
						PlanTime:     0,
						ScheduleTime: 0,
						StartTime:    0,
						EndTime:      0,
					}
				}

				log.Printf("exec task {%s} err: %v\n", res.Execution.Name, res.Error)
				// 此处不能return，出现异常时会导致任务中断
				continue
			}
			spent := res.EndTime.Sub(res.StartTime).Milliseconds()
			fmt.Printf("task {%v}, spent %v ms, output: %s\n", res.Execution.Name, spent, res.Output)
		}
		after = s.tryScheduler()
		// 重置timer
		timer.Reset(after)
	}
}

func (s *JobScheduler) tryScheduler() (after time.Duration) {
	if len(s.Plans) == 0 {
		return time.Second
	}
	now := time.Now()
	var nearestTime *time.Time
	for _, plan := range s.Plans {
		if plan.NextTime.Before(now) || plan.NextTime.Equal(now) {
			// @TODO 任务已到执行时间,尝试执行
			//fmt.Printf("run job: %s\n", plan.Job.Name)
			s.execute(plan)
			plan.NextTime = plan.Expression.Next(now)
		}
		// 获取下一个要执行的最近的任务时间
		if nearestTime == nil || plan.NextTime.Before(*nearestTime) {
			nearestTime = &plan.NextTime
		}
	}
	// 下次调度任务的时间间隔
	return (*nearestTime).Sub(now)
}

func (s *JobScheduler) eventHandler(e JobEvent) {
	switch e.Type {
	case PutEvent:
		plan, err := BuildJobPlan(e.Job)
		if err != nil {
			log.Println(err)
			return
		}
		s.Plans[e.Job.Name] = plan
	case DeleteEvent:
		if _, ok := s.Plans[e.Job.Name]; ok {
			delete(s.Plans, e.Job.Name)
		}
	case KillEvent:
		// 检查任务是否正在执行中
		if je, ok := s.Executions[e.Job.Name]; ok {
			je.CancelFunc()
		}
	}
}

func BuildJobPlan(job model.Job) (*JobPlan, error) {
	expr, err := cronexpr.Parse(job.Expression)
	if err != nil {
		return nil, fmt.Errorf("parse cron expr err: %v\n", err)
	}
	return &JobPlan{
		Job:        &job,
		Expression: expr,
		NextTime:   expr.Next(time.Now()),
	}, nil
}

func (s *JobScheduler) pushEvent(event JobEvent) {
	s.EventChan <- event
}
