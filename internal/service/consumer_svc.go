package service

import (
	"context"

	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle-layout/internal/tasks"
)

var Tasks map[string]*asynq.Task

var HandleFunc map[string]func(context.Context, *asynq.Task) error

type ConsumerService struct {
}

func NewConsumerService() *ConsumerService {
	return &ConsumerService{}
}

// Register register task and handle
// Tips: The number of tasks and handles must be the same
func (s *ConsumerService) Register() {
	// tasks
	Tasks = map[string]*asynq.Task{
		tasks.TypeEmailWelcome: tasks.NewEmailWelcomeTask(1),
	}

	// handles
	HandleFunc = map[string]func(context.Context, *asynq.Task) error{
		tasks.TypeEmailWelcome: tasks.HandleEmailWelcomeTask,
	}
}
