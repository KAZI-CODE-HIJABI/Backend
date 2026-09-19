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
	ChallengeID string `json:"challenge_id"`
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
}

type TestResults struct {
	Passed          int
	Failed          int
	Total           int
	ExecutionTimeMS int
}

func (r TestResults) Validate() error {
	if r.Passed < 0 || r.Failed < 0 || r.Total < 0 || r.ExecutionTimeMS < 0 {
		return errors.New("test results cannot contain negative values")
	}
	if r.Passed+r.Failed != r.Total {
		return errors.New("test result total must equal passed plus failed")
	}
	return nil
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

func CanClaimJob(status Status) bool  { return status == Queued }
func CanFinishJob(status Status) bool { return status == Running }
