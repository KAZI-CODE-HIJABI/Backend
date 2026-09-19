package assessment

import "testing"

func TestBuildUsesBackendRubricAndIntegrityRule(t *testing.T) {
	result, err := Build(Input{Rubric: RubricInput{CorrectImplementation: .9, TestsEdgeCases: .75, Debugging: .85, CodeQuality: .75, Explanation: .8}, HasRequiredEvidence: true, Feedback: "The core fix works; edge cases need attention.", Evidence: []Evidence{{Competency: "DEBUGGING", Status: Demonstrated, Explanation: "The explanation matches the tested fix."}}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Score != 82 || result.IntegrityStatus != IntegrityHigh {
		t.Fatalf("unexpected result: %#v", result)
	}
}
func TestBuildRequiresTraceableEvidence(t *testing.T) {
	if _, err := Build(Input{Rubric: RubricInput{}, Feedback: "feedback"}); err == nil {
		t.Fatal("assessment without evidence must be rejected")
	}
}
