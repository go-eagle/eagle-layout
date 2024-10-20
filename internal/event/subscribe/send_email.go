package subscribe

import (
	"context"
	"encoding/json"

	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

// 自定义消息处理函数
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
