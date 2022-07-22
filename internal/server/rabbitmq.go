package server

import (
	"context"
	"fmt"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

// NewRabbitMQServer creates a rabbitmq server
func NewRabbitMQServer() *rabbitmq.Server {
	srv := rabbitmq.NewServer(
		"guest:guest@localhost:5672",
		"test-exchange",
	)

	handler := func(ctx context.Context, body []byte) error {
		fmt.Println("handle msg: ", string(body))
		return nil
	}

	err := srv.RegisterSubscriber(context.Background(), "test-queue", handler)
	if err != nil {
		panic(err)
	}

	return srv
}
