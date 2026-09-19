// Package runner defines the trusted boundary around candidate-code execution.
// The API process must call an isolated implementation; it must never execute the candidate's source code itself.
package runner

import (
	"context"
	"errors"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
)

type Request struct {
	SubmissionID string
	ChallengeID  string
	Code         string
}
type IsolatedRunner interface {
	Run(context.Context, Request) (submission.TestResults, error)
}

func ValidateRequest(request Request) error {
	if request.SubmissionID == "" || request.ChallengeID == "" || request.Code == "" {
		return errors.New("runner request requires submission ID, challenge ID, and code")
	}
	return nil
}
