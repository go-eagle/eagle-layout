package subscribe

import (
	"context"
	"encoding/json"

	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

// SendWelcomeEmailHandler handles the welcome-email queue message.
func SendWelcomeEmailHandler(ctx context.Context, body amqp091.Delivery) (action rabbitmq.Action) {
	msg := make(map[string]interface{})
	err := json.Unmarshal(body.Body, &msg)
	if err != nil {
		logger.Errorf("consumer handler unmarshal msg err: %s", err.Error())
		return rabbitmq.NackDiscard
	}
	logger.Infof("consumer handler receive msg: %s", msg)

	// 下面可以增加自己的业务逻辑处理

	return rabbitmq.Ack
}
