package tasks

import (
	"context"
	"encoding/json"
	"fmt"

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

func NewEmailWelcomeTask(userID int) *asynq.Task {
	payload, _ := json.Marshal(EmailWelcomePayload{UserID: userID})
	return asynq.NewTask(TypeEmailWelcome, payload)
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
