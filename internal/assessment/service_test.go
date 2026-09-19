package assessment

import (
	"context"
	"testing"
)

type fakeRepository struct{ saved bool }

func (f *fakeRepository) Save(context.Context, string, string, Result, string, string) (StoredResult, error) {
	f.saved = true
	return StoredResult{AssessmentID: "assessment"}, nil
}
func (f *fakeRepository) Get(context.Context, string, string) (StoredResult, error) {
	return StoredResult{AssessmentID: "assessment"}, nil
}
func TestServiceBuildsBeforeSaving(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo)
	_, err := service.Save(context.Background(), "candidate", "submission", Input{Feedback: "feedback", Evidence: []Evidence{{Competency: "DEBUGGING", Status: Demonstrated, Explanation: "evidence"}}}, "v1", "")
	if err != nil {
		t.Fatal(err)
	}
	if !repo.saved {
		t.Fatal("valid assessment must be saved")
	}
}
