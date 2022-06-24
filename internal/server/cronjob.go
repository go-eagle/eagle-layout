package server

import (
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/cronjob"
	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/tasks"
)

type Server struct {
	c    *tasks.Config
	sche *asynq.Scheduler
}

func NewCronJobServer(c *tasks.Config, jobSvc *service.CronJobService) (s *cronjob.Server) {
	jobSvc.RegisterTask()
	s = cronjob.NewServer()
	for _, j := range c.Jobs {
		job, ok := service.DefaultJobs[j.Name]
		if !ok {
			log.Warnf("can not find job: %s", j.Name)
			continue
		}
		_, err := s.RegisterJob(j.Schedule, job)
		if err != nil {
			panic(err)
		}
	}
	log.Info("load jobs count:", len(c.Jobs))
	return s
}
