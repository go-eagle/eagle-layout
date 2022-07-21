package server

import (
	"time"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/cronjob"
	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/tasks"
)

type Server struct {
	c      *tasks.Config
	sche   *asynq.Scheduler
	client *asynq.Client
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

	srv := &Server{
		c:      c,
		client: tasks.NewClient(c),
	}
	srv.RunWorkerServer(c)
	srv.RunDelayTasks()

	return s
}

// RunWorkerServer Run worker server
func (s *Server) RunWorkerServer(c *tasks.Config) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Addr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: c.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				tasks.QueueCritical: 6,
				tasks.QueueDefault:  3,
				tasks.QueueLow:      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	// register handlers...
	mux.HandleFunc(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Errorf("could not run async server: %v", err)
		}
	}()
}

// RunDelayTasks Run delay tasks
func (s *Server) RunDelayTasks() {
	// ------------------------------------------------------
	// Enqueue task to be processed immediately.
	// Use (*Client).Enqueue method.
	// ------------------------------------------------------
	task := tasks.NewEmailWelcomeTask(1)
	info, err := s.client.Enqueue(task, asynq.Queue(tasks.QueueDefault))
	if err != nil {
		log.Errorf("could not enqueue task: %v", err)
	}
	log.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Schedule task to be processed in the future.
	// Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	task = tasks.NewEmailWelcomeTask(2)
	info, err = s.client.Enqueue(task, asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Errorf("could not enqueue task: %v", err)
	}
	log.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Set other options to tune task processing behavior.
	// Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------
	task = tasks.NewEmailWelcomeTask(3)
	info, err = s.client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Errorf("could not enqueue task: %v", err)
	}
	log.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
