package server

import (
	"github.com/go-eagle/eagle-layout/internal/tasks"
	logger "github.com/go-eagle/eagle/pkg/log"
	redisMQ "github.com/go-eagle/eagle/pkg/transport/consumer/redis"
	"github.com/hibiken/asynq"
)

// NewRedisConsumerServer create a redis server
func NewRedisConsumerServer(c *tasks.Config) *redisMQ.Server {
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
			// Logger: logger.GetLogger(),
			// See the godoc for other configuration options
		},
	)

	// register task
	// 创建任务
	err := tasks.NewEmailWelcomeTask(tasks.EmailWelcomePayload{UserID: 1})
	if err != nil {
		logger.Fatalf("could not create task: %v", err)
	}

	// register handler
	srv.RegisterHandler(tasks.TypeEmailWelcome, tasks.HandleEmailWelcomeTask)
	// here register other handlers...

	return srv
}
