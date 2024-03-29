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
	TypeEmailWelcome = "email:welcome"
)

type EmailWelcomePayload struct {
	UserID int
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewEmailWelcomeTask(data EmailWelcomePayload) error {
	payload, _ := json.Marshal(data)

	task := asynq.NewTask(TypeEmailWelcome, payload)
	_, err := GetClient().Enqueue(task)
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

func HandleEmailWelcomeTask(ctx context.Context, t *asynq.Task) error {
	var p EmailWelcomePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Infof("Sending Email to User: user_id=%d", p.UserID)
	// Email delivery code ...
	return nil
}
