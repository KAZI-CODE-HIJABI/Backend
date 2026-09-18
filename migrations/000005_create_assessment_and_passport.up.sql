CREATE TABLE assessments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  submission_id UUID NOT NULL UNIQUE REFERENCES submissions(id) ON DELETE CASCADE,
  score INTEGER NOT NULL CHECK (score BETWEEN 0 AND 100),
  integrity_status TEXT NOT NULL CHECK (integrity_status IN ('HIGH', 'REVIEW_RECOMMENDED')),
  feedback TEXT NOT NULL,
  rules_version TEXT NOT NULL,
  prompt_version TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE assessment_evidence (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
  competency TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('STRONG_EVIDENCE', 'DEMONSTRATED', 'DEVELOPING', 'NEEDS_IMPROVEMENT')),
  explanation TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (assessment_id, competency)
);

CREATE TABLE skills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  category TEXT NOT NULL DEFAULT 'SOFTWARE_ENGINEERING',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE candidate_skills (
  candidate_id UUID NOT NULL REFERENCES candidate_profiles(candidate_id) ON DELETE CASCADE,
  skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE RESTRICT,
  assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('STRONG_EVIDENCE', 'DEMONSTRATED', 'DEVELOPING', 'NEEDS_IMPROVEMENT')),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (candidate_id, skill_id)
);

CREATE INDEX assessment_evidence_assessment_id_idx ON assessment_evidence (assessment_id);
CREATE INDEX candidate_skills_candidate_id_idx ON candidate_skills (candidate_id);
