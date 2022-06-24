package service

import (
	"fmt"

	"github.com/go-eagle/eagle-layout/internal/tasks"

	"github.com/hibiken/asynq"
)

var DefaultJobs map[string]*asynq.Task

type JobFunc func()

type CronJobService struct {
}

func NewCronJobService() *CronJobService {
	return &CronJobService{}
}

func (s *CronJobService) RegisterTask() {
	DefaultJobs = map[string]*asynq.Task{
		"demo": tasks.NewDemoJobTask(),
	}
}

func (s *CronJobService) Demo() {
	fmt.Println("this is a demo")
}
