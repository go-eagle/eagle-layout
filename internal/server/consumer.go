package server

import (
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/consumer"
	"github.com/hibiken/asynq"
)

func NewConsumerServer(c *tasks.Config, jobSvc *service.ConsumerService) *consumer.Server {
	// register task and handle
	jobSvc.Register()

	cfg := asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    c.Redis.Addr,
	}

	srv := consumer.NewServer(cfg, asynq.Config{
		// Specify how many concurrent workers to use
		Concurrency: c.Redis.Concurrency,
		// Optionally specify multiple queues with different priority.
		Queues: map[string]int{
			tasks.QueueCritical: 6,
			tasks.QueueDefault:  3,
			tasks.QueueLow:      1,
		},
		// See the godoc for other configuration options
	})
	for _, t := range c.Tasks {
		task, ok := service.Tasks[t.Name]
		if !ok {
			log.Warnf("can not find task: %s", t.Name)
			continue
		}
		// register task
		_, err := srv.RegisterTask(t.Schedule, task)
		if err != nil {
			panic(err)
		}

		// register handle
		handle, ok := service.HandleFunc[t.Name]
		if !ok {
			log.Warnf("can not find handle: %s", t.Name)
			continue
		}
		log.Info("====register handle===", t.Name)
		srv.RegisterHandle(t.Name, handle)
	}
	log.Info("load tasks count:", len(c.Tasks))

	//go func() {
	//	s := asynq.NewServer(
	//		cfg,
	//		asynq.Config{
	//
	//		},
	//	)
	//
	//	// mux maps a type to a handler
	//	mux := asynq.NewServeMux()
	//	// register handlers...
	//	mux.HandleFunc(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)
	//
	//	if err := s.Run(mux); err != nil {
	//		log.Errorf("failed to run async server: %v", err)
	//	}
	//}()

	return srv
}
