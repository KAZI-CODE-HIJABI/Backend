package submission

import (
	"context"
	"errors"
)

var (
	ErrNotFound            = errors.New("submission not found")
	ErrActiveSubmission    = errors.New("candidate already has an active submission for this challenge")
	ErrFollowupUnavailable = errors.New("follow-up is unavailable")
)

type View struct {
	SubmissionID   string
	ChallengeID    string
	Status         Status
	TestResults    *TestResults
	FollowupStatus string
}

type Repository interface {
	Create(context.Context, string, CreateInput) (View, error)
	Get(context.Context, string, string) (View, error)
	SaveFollowup(context.Context, string, string, string) error
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Create(ctx context.Context, candidateID string, input CreateInput) (View, error) {
	if candidateID == "" {
		return View{}, errors.New("candidate identity is required")
	}
	if err := input.Validate(); err != nil {
		return View{}, err
	}
	return s.repository.Create(ctx, candidateID, input)
}
func (s *Service) Get(ctx context.Context, candidateID, submissionID string) (View, error) {
	if candidateID == "" || submissionID == "" {
		return View{}, ErrNotFound
	}
	return s.repository.Get(ctx, candidateID, submissionID)
}
func (s *Service) AnswerFollowup(ctx context.Context, candidateID, submissionID, answer string) error {
	if candidateID == "" || submissionID == "" || answer == "" {
		return ErrFollowupUnavailable
	}
	return s.repository.SaveFollowup(ctx, candidateID, submissionID, answer)
}
