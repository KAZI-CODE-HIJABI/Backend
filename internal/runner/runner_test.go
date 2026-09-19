package runner

import "testing"

func TestValidateRequest(t *testing.T) {
	if err := ValidateRequest(Request{SubmissionID: "submission", ChallengeID: "challenge", Code: "module.exports = {}"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateRequest(Request{}); err == nil {
		t.Fatal("empty execution request must be rejected")
	}
}
