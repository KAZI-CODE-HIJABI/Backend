package submission

import "testing"

func TestCreateInputValidation(t *testing.T) {
	if err := (CreateInput{ChallengeID: "challenge", Code: "console.log(1)", Explanation: "reason"}).Validate(); err != nil {
		t.Fatal(err)
	}
	for _, input := range []CreateInput{{}, {ChallengeID: "c"}, {ChallengeID: "c", Code: "x"}} {
		if input.Validate() == nil {
			t.Fatal("expected validation error")
		}
	}
}

func TestTransitions(t *testing.T) {
	if err := Transition(Queued, Running); err != nil {
		t.Fatal(err)
	}
	if err := Transition(Running, Completed); err != nil {
		t.Fatal(err)
	}
	if err := Transition(Completed, Running); err == nil {
		t.Fatal("completed submission must be terminal")
	}
}

func TestFollowupRules(t *testing.T) {
	if err := CanAnswerFollowup(Completed, false); err != nil {
		t.Fatal(err)
	}
	if err := CanAnswerFollowup(Running, false); err == nil {
		t.Fatal("running submission must not accept follow-up")
	}
	if err := CanAnswerFollowup(Completed, true); err == nil {
		t.Fatal("repeat answer must be rejected")
	}
}

func TestTestResultsMustBeInternallyConsistent(t *testing.T) {
	if err := (TestResults{Passed: 8, Failed: 2, Total: 10, ExecutionTimeMS: 1820}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (TestResults{Passed: 8, Failed: 2, Total: 9}).Validate(); err == nil {
		t.Fatal("mismatched result total must be rejected")
	}
}
