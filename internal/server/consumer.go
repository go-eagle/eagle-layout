package server

import (
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/consumer"
	"github.com/hibiken/asynq"
)

func NewConsumerServer(c *tasks.Config, svc *service.ConsumerService) *consumer.Server {
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

	// register handle
	handles := svc.RegisterHandle()

	for name, handle := range handles {
		log.Info("[server] consumer register handle: ", name)
		srv.RegisterHandle(name, handle)
	}
	log.Info("[server] consumer load tasks count:", len(handles))

	return srv
}
