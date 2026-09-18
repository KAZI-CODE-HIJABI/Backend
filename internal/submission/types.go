package submission

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	Queued    Status = "QUEUED"
	Running   Status = "RUNNING"
	Completed Status = "COMPLETED"
	Failed    Status = "FAILED"
)

type CreateInput struct {
	ChallengeID string
	Code        string
	Explanation string
}

func (in CreateInput) Validate() error {
	if strings.TrimSpace(in.ChallengeID) == "" {
		return errors.New("challenge_id is required")
	}
	if strings.TrimSpace(in.Code) == "" {
		return errors.New("code is required")
	}
	if strings.TrimSpace(in.Explanation) == "" {
		return errors.New("explanation is required")
	}
	return nil
}

func CanTransition(from, to Status) bool {
	return (from == Queued && (to == Running || to == Failed)) ||
		(from == Running && (to == Completed || to == Failed))
}

func Transition(from, to Status) error {
	if !CanTransition(from, to) {
		return fmt.Errorf("invalid submission status transition: %s to %s", from, to)
	}
	return nil
}

func CanAnswerFollowup(status Status, existingAnswer bool) error {
	if status != Completed {
		return errors.New("submission execution is not completed")
	}
	if existingAnswer {
		return errors.New("follow-up was already answered")
	}
	return nil
}
