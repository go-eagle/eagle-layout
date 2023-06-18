package service

import (
	"context"
	"fmt"

	"github.com/go-eagle/eagle/pkg/transport/consumer"

	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle-layout/internal/tasks"
)

type ConsumerService struct {
}

func NewConsumerService() *ConsumerService {
	return &ConsumerService{}
}

// RegisterHandle register handle
func (s *ConsumerService) RegisterHandle() map[string]func(context.Context, *asynq.Task) error {
	// config handles
	HandleFunc := map[string]func(context.Context, *asynq.Task) error{
		tasks.TypeEmailWelcome: tasks.HandleEmailWelcomeTask,
	}

	return HandleFunc
}

// RegisterSchedule register schedule task
func (s *ConsumerService) RegisterSchedule(srv *consumer.Server, name string, task *asynq.Task) (string, error) {
	// config schedules
	schedules := map[string]string{
		tasks.TypeEmailWelcome: "@every 10s",
	}

	schedule, ok := schedules[name]
	if !ok {
		return "", fmt.Errorf("no register schedule, name: %s", name)
	}

	return srv.RegisterTask(schedule, task)
}
