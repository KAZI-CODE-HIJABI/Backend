package passport

import "github.com/KAZI-CODE-HIJABI/Backend/internal/assessment"

type Skill struct {
	Name   string                    `json:"name"`
	Status assessment.EvidenceStatus `json:"status"`
}
type Summary struct {
	ChallengesCompleted  int    `json:"challenges_completed"`
	TestsPassed          string `json:"tests_passed"`
	ExplanationSubmitted bool   `json:"explanation_submitted"`
	FollowupCompleted    bool   `json:"followup_completed"`
}
type View struct {
	CandidateID     string  `json:"candidate_id"`
	AnonymousID     string  `json:"anonymous_id"`
	Category        string  `json:"category"`
	OverallScore    int     `json:"overall_score"`
	Skills          []Skill `json:"skills"`
	EvidenceSummary Summary `json:"evidence_summary"`
}

// NewView projects completed assessment data for the candidate's own passport.
func NewView(candidateID, anonymousID string, score int, skills []Skill, summary Summary) View {
	return View{CandidateID: candidateID, AnonymousID: anonymousID, Category: "SOFTWARE_ENGINEERING", OverallScore: score, Skills: skills, EvidenceSummary: summary}
}
