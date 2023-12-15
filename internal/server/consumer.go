package server

import (
	"github.com/go-eagle/eagle-layout/internal/tasks"
	redisMQ "github.com/go-eagle/eagle/pkg/transport/consumer/redis"
	"github.com/hibiken/asynq"
)

func NewConsumerServer(c *tasks.Config) *redisMQ.Server {
	cfg := asynq.RedisClientOpt{
		Network: "tcp",
		Addr:    c.Redis.Addr,
	}

	srv := redisMQ.NewServer(
		cfg,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: c.Redis.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				redisMQ.QueueCritical: 6,
				redisMQ.QueueDefault:  3,
				redisMQ.QueueLow:      1,
			},
			// See the godoc for other configuration options
		},
	)

	// register handler
	srv.RegisterHandler(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)
	// here register other handlers...

	return srv
}
