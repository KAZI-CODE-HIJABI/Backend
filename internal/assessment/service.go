package assessment

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("assessment not found")

type StoredResult struct {
	AssessmentID string
	Result       Result
}
type Repository interface {
	Save(context.Context, string, string, Result, string, string) (StoredResult, error)
	Get(context.Context, string, string) (StoredResult, error)
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Save(ctx context.Context, candidateID, submissionID string, input Input, rulesVersion, promptVersion string) (StoredResult, error) {
	if candidateID == "" || submissionID == "" || rulesVersion == "" {
		return StoredResult{}, errors.New("candidate, submission, and rules version are required")
	}
	result, err := Build(input)
	if err != nil {
		return StoredResult{}, err
	}
	return s.repository.Save(ctx, candidateID, submissionID, result, rulesVersion, promptVersion)
}
func (s *Service) Get(ctx context.Context, candidateID, assessmentID string) (StoredResult, error) {
	if candidateID == "" || assessmentID == "" {
		return StoredResult{}, ErrNotFound
	}
	return s.repository.Get(ctx, candidateID, assessmentID)
}
