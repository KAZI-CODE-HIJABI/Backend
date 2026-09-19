package ai

import (
	"context"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
	"testing"
)

type fakeProvider struct {
	observations []Observation
	feedback     string
}

func (f fakeProvider) Analyze(context.Context, AnalysisRequest) ([]Observation, string, error) {
	return f.observations, f.feedback, nil
}
func completeRequest() AnalysisRequest {
	return AnalysisRequest{ChallengeDescription: "Fix the API", RubricVersion: "v1", Code: "code", Explanation: "reason", FollowupQuestion: "what if missing", FollowupAnswer: "return 404", TestResults: submission.TestResults{Passed: 8, Failed: 2, Total: 10, ExecutionTimeMS: 1820}}
}
func TestAnalyzeOnlyReturnsValidatedEvidence(t *testing.T) {
	service := NewService(fakeProvider{observations: []Observation{{Competency: "DEBUGGING", Finding: "The explanation matches the fixed lookup.", Support: []string{"code:1"}}}, feedback: "Good fix."})
	observations, feedback, err := service.Analyze(context.Background(), completeRequest())
	if err != nil {
		t.Fatal(err)
	}
	if len(observations) != 1 || feedback != "Good fix." {
		t.Fatal("unexpected AI result")
	}
}
func TestAnalyzeRejectsIncompleteEvidence(t *testing.T) {
	service := NewService(fakeProvider{})
	request := completeRequest()
	request.FollowupAnswer = ""
	if _, _, err := service.Analyze(context.Background(), request); err == nil {
		t.Fatal("incomplete assessment must not reach AI")
	}
}
