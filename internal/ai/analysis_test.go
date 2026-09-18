package ai

import "testing"

func TestValidateObservations(t *testing.T) {
	valid := []Observation{{Competency: "DEBUGGING", Finding: "The answer explains the failed lookup.", Support: []string{"explanation:1"}}}
	if err := Validate(valid); err != nil {
		t.Fatal(err)
	}
	if err := Validate([]Observation{{Competency: "FINAL_SCORE", Finding: "82", Support: []string{"ai"}}}); err == nil {
		t.Fatal("must reject score request")
	}
	if err := Validate([]Observation{{Competency: "DEBUGGING", Finding: "x"}}); err == nil {
		t.Fatal("must require support")
	}
}
