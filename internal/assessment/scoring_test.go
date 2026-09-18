package assessment

import "testing"

func TestScoreUsesFixedRubric(t *testing.T) {
	score, err := (RubricInput{CorrectImplementation: .9, TestsEdgeCases: .75, Debugging: .85, CodeQuality: .75, Explanation: .8}).Score()
	if err != nil {
		t.Fatal(err)
	}
	if score != 82 {
		t.Fatalf("got %d, want 82", score)
	}
}
func TestScoreRejectsOutOfRangeInput(t *testing.T) {
	if _, err := (RubricInput{CorrectImplementation: 1.1}).Score(); err == nil {
		t.Fatal("expected validation error")
	}
}
func TestEvidenceStatuses(t *testing.T) {
	cases := []struct {
		score float64
		want  EvidenceStatus
	}{{.9, StrongEvidence}, {.6, Demonstrated}, {.35, Developing}, {.1, NeedsImprovement}}
	for _, tc := range cases {
		got, err := StatusFor(tc.score)
		if err != nil || got != tc.want {
			t.Fatalf("%v: got %q, %v", tc.score, got, err)
		}
	}
}
func TestIntegrityNeverUsesCheatingClaim(t *testing.T) {
	if Integrity(true, false) != IntegrityHigh {
		t.Fatal("complete consistent evidence should be high")
	}
	if Integrity(false, false) != ReviewRecommended {
		t.Fatal("missing evidence should recommend review")
	}
	if Integrity(true, true) != ReviewRecommended {
		t.Fatal("contradiction should recommend review")
	}
}
