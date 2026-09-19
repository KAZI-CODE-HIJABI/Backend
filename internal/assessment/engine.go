package assessment

import "errors"

type Evidence struct {
	Competency  string         `json:"competency"`
	Status      EvidenceStatus `json:"status"`
	Explanation string         `json:"explanation"`
}
type Input struct {
	Rubric                   RubricInput
	HasRequiredEvidence      bool
	HasMaterialContradiction bool
	Feedback                 string
	Evidence                 []Evidence
}
type Result struct {
	Score           int             `json:"score"`
	IntegrityStatus IntegrityStatus `json:"integrity_status"`
	Feedback        string          `json:"feedback"`
	Evidence        []Evidence      `json:"evidence"`
}

func Build(input Input) (Result, error) {
	if input.Feedback == "" {
		return Result{}, errors.New("assessment feedback is required")
	}
	if len(input.Evidence) == 0 {
		return Result{}, errors.New("assessment evidence is required")
	}
	seen := make(map[string]struct{}, len(input.Evidence))
	for _, e := range input.Evidence {
		if e.Competency == "" || e.Explanation == "" {
			return Result{}, errors.New("assessment evidence requires competency and explanation")
		}
		if _, exists := seen[e.Competency]; exists {
			return Result{}, errors.New("assessment cannot contain duplicate competency evidence")
		}
		seen[e.Competency] = struct{}{}
	}
	score, err := input.Rubric.Score()
	if err != nil {
		return Result{}, err
	}
	return Result{Score: score, IntegrityStatus: Integrity(input.HasRequiredEvidence, input.HasMaterialContradiction), Feedback: input.Feedback, Evidence: input.Evidence}, nil
}
