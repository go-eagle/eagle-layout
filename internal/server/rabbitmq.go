package server

import (
	"fmt"

	rabbitmqConf "github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/transport/consumer/rabbitmq"

	"github.com/go-eagle/eagle-layout/internal/event/subscribe"
	"github.com/go-eagle/eagle-layout/internal/tasks"
)

// NewRabbitmqConsumerServer create a redis server
func NewRabbitmqConsumerServer() *rabbitmq.Server {
	rabbitmqConf.Load()

	srv := rabbitmq.NewServer()

	if err := srv.RegisterHandler(tasks.TypeEmailWelcome, subscribe.SendWelcomeEmailHandler); err != nil {
		panic(fmt.Errorf("register welcome email handler: %w", err))
	}
	// here register other handlers...

	return srv
}
