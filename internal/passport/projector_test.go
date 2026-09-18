package passport

import (
	"github.com/KAZI-CODE-HIJABI/Backend/internal/assessment"
	"testing"
)

func TestNewView(t *testing.T) {
	v := NewView("candidate", "KZ-1042", 82, []Skill{{Name: "Debugging", Status: assessment.Demonstrated}}, Summary{ChallengesCompleted: 1, TestsPassed: "8/10", ExplanationSubmitted: true, FollowupCompleted: true})
	if v.Category != "SOFTWARE_ENGINEERING" || v.OverallScore != 82 || v.AnonymousID != "KZ-1042" {
		t.Fatalf("unexpected passport: %#v", v)
	}
}
