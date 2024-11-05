package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/hibiken/asynq"
)

const (
	// TypeEmailWelcome define a task name
	TypeEmailWelcome = "email:welcome"
)

// EmailWelcomePayload payload data for task
type EmailWelcomePayload struct {
	UserID int64
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

// NewEmailWelcomeTask create a task
func NewEmailWelcomeTask(data EmailWelcomePayload) error {
	payload, _ := json.Marshal(data)

	task := asynq.NewTask(TypeEmailWelcome, payload)
	_, err := GetClient().Enqueue(task)

	// // 即时消息
	// _, err := GetClient().Enqueue(task)

	// // 延时消息
	// _, err := GetClient().Enqueue(task, asynq.ProcessIn(10*time.Second))

	// // 定时消息
	// _, err := GetClient().Enqueue(task, asynq.ProcessAt(time.Now().Add(time.Hour)))

	// // 超时、重试
	// _, err := GetClient().Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))

	// // 优先级消息
	// _, err := GetClient().Enqueue(task, asynq.Queue(QueueCritical))

	if err != nil {
		return errors.Wrapf(err, "[tasks] Enqueue task error, name: %s", TypeEmailWelcome)
	}

	return err
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

// HandleEmailWelcomeTask handle task
func HandleEmailWelcomeTask(ctx context.Context, t *asynq.Task) error {
	var p EmailWelcomePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Infof("Sending Email to User: user_id=%d", p.UserID)
	// Email delivery code ...
	return nil
}
