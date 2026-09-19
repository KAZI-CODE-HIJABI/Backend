package ai

import (
	"context"
	"errors"

	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
)

// AnalysisRequest deliberately excludes candidate identity and employer data.
type AnalysisRequest struct {
	ChallengeDescription string
	RubricVersion        string
	Code                 string
	Explanation          string
	FollowupQuestion     string
	FollowupAnswer       string
	TestResults          submission.TestResults
}

type Provider interface {
	Analyze(context.Context, AnalysisRequest) ([]Observation, string, error)
}
type Service struct{ provider Provider }

func NewService(provider Provider) *Service { return &Service{provider: provider} }
func (s *Service) Analyze(ctx context.Context, request AnalysisRequest) ([]Observation, string, error) {
	if s.provider == nil {
		return nil, "", errors.New("AI provider is not configured")
	}
	if request.ChallengeDescription == "" || request.RubricVersion == "" || request.Code == "" || request.Explanation == "" || request.FollowupQuestion == "" || request.FollowupAnswer == "" {
		return nil, "", errors.New("AI analysis requires the complete assessment evidence set")
	}
	if err := request.TestResults.Validate(); err != nil {
		return nil, "", err
	}
	observations, feedback, err := s.provider.Analyze(ctx, request)
	if err != nil {
		return nil, "", err
	}
	if err := Validate(observations); err != nil {
		return nil, "", err
	}
	if feedback == "" {
		return nil, "", errors.New("AI feedback is required")
	}
	return observations, feedback, nil
}
