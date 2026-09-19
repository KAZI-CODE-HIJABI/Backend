package passport

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("skills passport not found")

type Repository interface {
	Get(context.Context, string) (View, error)
}
type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }
func (s *Service) Get(ctx context.Context, candidateID string) (View, error) {
	if candidateID == "" {
		return View{}, ErrNotFound
	}
	return s.repository.Get(ctx, candidateID)
}
