package submission

import (
	"context"
	"testing"
)

type fakeRepository struct{ created bool }

func (f *fakeRepository) Create(context.Context, string, CreateInput) (View, error) {
	f.created = true
	return View{SubmissionID: "submission", Status: Queued}, nil
}
func (f *fakeRepository) Get(context.Context, string, string) (View, error) {
	return View{SubmissionID: "submission", Status: Completed}, nil
}
func (f *fakeRepository) SaveFollowup(context.Context, string, string, string) error { return nil }
func TestServiceValidatesBeforePersistence(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)
	if _, err := service.Create(context.Background(), "candidate", CreateInput{}); err == nil {
		t.Fatal("invalid input must be rejected")
	}
	if repo.created {
		t.Fatal("invalid input must not be persisted")
	}
	if _, err := service.Create(context.Background(), "candidate", CreateInput{ChallengeID: "challenge", Code: "code", Explanation: "explanation"}); err != nil {
		t.Fatal(err)
	}
}
