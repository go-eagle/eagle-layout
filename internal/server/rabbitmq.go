package server

import (
	"github.com/go-eagle/eagle-layout/internal/event/subscribe"
	"github.com/go-eagle/eagle-layout/internal/tasks"
	rabbitmqConf "github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/transport/consumer/rabbitmq"
)

// NewRabbitmqConsumerServer create a redis server
func NewRabbitmqConsumerServer() *rabbitmq.Server {
	rabbitmqConf.Load()

	srv := rabbitmq.NewServer()

	// register handler
	srv.RegisterHandler(tasks.TypeEmailWelcome, subscribe.SendWelcomeEmailHandler)
	// here register other handlers...

	return srv
}
