package ai

import (
	"errors"
	"strings"
)

// Observation is qualitative input to KAZI's rules, never a final grade, rank,
// hiring decision, cheating probability, or integrity decision.
type Observation struct {
	Competency string   `json:"competency"`
	Finding    string   `json:"finding"`
	Support    []string `json:"supporting_references"`
}

var allowedCompetencies = map[string]struct{}{
	"CORRECT_IMPLEMENTATION": {}, "TESTS_EDGE_CASES": {}, "DEBUGGING": {}, "CODE_QUALITY": {}, "EXPLANATION": {},
}

func Validate(observations []Observation) error {
	seen := make(map[string]struct{})
	for _, observation := range observations {
		if _, ok := allowedCompetencies[observation.Competency]; !ok {
			return errors.New("unsupported competency returned by AI")
		}
		if strings.TrimSpace(observation.Finding) == "" {
			return errors.New("AI observation is missing a finding")
		}
		if len(observation.Support) == 0 {
			return errors.New("AI observation must cite submitted or trusted evidence")
		}
		if _, exists := seen[observation.Competency]; exists {
			return errors.New("AI returned duplicate competency observation")
		}
		seen[observation.Competency] = struct{}{}
	}
	return nil
}
