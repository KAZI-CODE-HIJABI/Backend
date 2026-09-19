package assessment

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type PGRepository struct{ pool *pgxpool.Pool }

func NewPGRepository(pool *pgxpool.Pool) *PGRepository { return &PGRepository{pool: pool} }
func (r *PGRepository) Save(ctx context.Context, candidateID, submissionID string, result Result, rulesVersion, promptVersion string) (StoredResult, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return StoredResult{}, err
	}
	defer tx.Rollback(ctx)
	var assessmentID string
	err = tx.QueryRow(ctx, `INSERT INTO assessments (submission_id,score,integrity_status,feedback,rules_version,prompt_version) SELECT s.id,$3,$4,$5,$6,NULLIF($7,'') FROM submissions s WHERE s.id=$1 AND s.candidate_id=$2 RETURNING id::text`, submissionID, candidateID, result.Score, result.IntegrityStatus, result.Feedback, rulesVersion, promptVersion).Scan(&assessmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return StoredResult{}, ErrNotFound
	}
	if err != nil {
		return StoredResult{}, fmt.Errorf("create assessment: %w", err)
	}
	for _, evidence := range result.Evidence {
		if _, err = tx.Exec(ctx, `INSERT INTO assessment_evidence (assessment_id,competency,status,explanation) VALUES ($1,$2,$3,$4)`, assessmentID, evidence.Competency, evidence.Status, evidence.Explanation); err != nil {
			return StoredResult{}, fmt.Errorf("save assessment evidence: %w", err)
		}
		skillName := skillName(evidence.Competency)
		var skillID string
		if err = tx.QueryRow(ctx, `INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id::text`, skillName).Scan(&skillID); err != nil {
			return StoredResult{}, fmt.Errorf("upsert skill: %w", err)
		}
		if _, err = tx.Exec(ctx, `INSERT INTO candidate_skills (candidate_id,skill_id,assessment_id,status) VALUES ($1,$2,$3,$4) ON CONFLICT (candidate_id,skill_id) DO UPDATE SET assessment_id=EXCLUDED.assessment_id,status=EXCLUDED.status,updated_at=now()`, candidateID, skillID, assessmentID, evidence.Status); err != nil {
			return StoredResult{}, fmt.Errorf("upsert candidate skill: %w", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return StoredResult{}, err
	}
	return StoredResult{AssessmentID: assessmentID, Result: result}, nil
}
func (r *PGRepository) Get(ctx context.Context, candidateID, assessmentID string) (StoredResult, error) {
	var stored StoredResult
	err := r.pool.QueryRow(ctx, `SELECT a.id::text,a.score,a.integrity_status,a.feedback FROM assessments a JOIN submissions s ON s.id=a.submission_id WHERE a.id=$1 AND s.candidate_id=$2`, assessmentID, candidateID).Scan(&stored.AssessmentID, &stored.Result.Score, &stored.Result.IntegrityStatus, &stored.Result.Feedback)
	if errors.Is(err, pgx.ErrNoRows) {
		return StoredResult{}, ErrNotFound
	}
	if err != nil {
		return StoredResult{}, fmt.Errorf("get assessment: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT competency,status,explanation FROM assessment_evidence WHERE assessment_id=$1 ORDER BY competency`, assessmentID)
	if err != nil {
		return StoredResult{}, fmt.Errorf("get assessment evidence: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var evidence Evidence
		if err := rows.Scan(&evidence.Competency, &evidence.Status, &evidence.Explanation); err != nil {
			return StoredResult{}, err
		}
		stored.Result.Evidence = append(stored.Result.Evidence, evidence)
	}
	if err := rows.Err(); err != nil {
		return StoredResult{}, err
	}
	return stored, nil
}
func skillName(competency string) string {
	switch competency {
	case "CORRECT_IMPLEMENTATION":
		return "API Development"
	case "TESTS_EDGE_CASES":
		return "Testing"
	case "CODE_QUALITY":
		return "Code Quality"
	default:
		return strings.Title(strings.ToLower(strings.ReplaceAll(competency, "_", " ")))
	}
}
