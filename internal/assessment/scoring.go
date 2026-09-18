package assessment

import (
	"errors"
	"fmt"
	"math"
)

type EvidenceStatus string
type IntegrityStatus string

const (
	StrongEvidence    EvidenceStatus  = "STRONG_EVIDENCE"
	Demonstrated      EvidenceStatus  = "DEMONSTRATED"
	Developing        EvidenceStatus  = "DEVELOPING"
	NeedsImprovement  EvidenceStatus  = "NEEDS_IMPROVEMENT"
	IntegrityHigh     IntegrityStatus = "HIGH"
	ReviewRecommended IntegrityStatus = "REVIEW_RECOMMENDED"
)

type RubricInput struct {
	CorrectImplementation float64
	TestsEdgeCases        float64
	Debugging             float64
	CodeQuality           float64
	Explanation           float64
}

// Score is KAZI's deterministic rule. AI observations may support a component,
// but callers must never pass a model-produced final score.
func (r RubricInput) Score() (int, error) {
	for name, value := range map[string]float64{
		"correct implementation": r.CorrectImplementation, "tests and edge cases": r.TestsEdgeCases,
		"debugging": r.Debugging, "code quality": r.CodeQuality, "explanation": r.Explanation,
	} {
		if value < 0 || value > 1 {
			return 0, fmt.Errorf("%s score must be between 0 and 1", name)
		}
	}
	total := r.CorrectImplementation*.30 + r.TestsEdgeCases*.20 + r.Debugging*.20 + r.CodeQuality*.15 + r.Explanation*.15
	return int(math.Round(total * 100)), nil
}

func StatusFor(score float64) (EvidenceStatus, error) {
	if score < 0 || score > 1 {
		return "", errors.New("evidence score must be between 0 and 1")
	}
	switch {
	case score >= .85:
		return StrongEvidence, nil
	case score >= .60:
		return Demonstrated, nil
	case score >= .35:
		return Developing, nil
	default:
		return NeedsImprovement, nil
	}
}

// Integrity is a transparent rule based on evidence consistency. It never claims
// to determine whether a candidate used AI or otherwise cheated.
func Integrity(hasRequiredEvidence, hasMaterialContradiction bool) IntegrityStatus {
	if !hasRequiredEvidence || hasMaterialContradiction {
		return ReviewRecommended
	}
	return IntegrityHigh
}
