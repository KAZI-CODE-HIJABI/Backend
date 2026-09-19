package submission

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepository struct{ pool *pgxpool.Pool }

func NewPGRepository(pool *pgxpool.Pool) *PGRepository { return &PGRepository{pool: pool} }

func (r *PGRepository) Create(ctx context.Context, candidateID string, input CreateInput) (View, error) {
	var view View
	err := r.pool.QueryRow(ctx, `INSERT INTO submissions (candidate_id, challenge_id, code, explanation) VALUES ($1,$2,$3,$4) RETURNING id::text, challenge_id::text, status`, candidateID, input.ChallengeID, input.Code, input.Explanation).Scan(&view.SubmissionID, &view.ChallengeID, &view.Status)
	if uniqueViolation(err) {
		return View{}, ErrActiveSubmission
	}
	if err != nil {
		return View{}, fmt.Errorf("create submission: %w", err)
	}
	view.FollowupStatus = "PENDING"
	return view, nil
}

func (r *PGRepository) Get(ctx context.Context, candidateID, submissionID string) (View, error) {
	var view View
	var passed, failed, total, executionMS *int
	err := r.pool.QueryRow(ctx, `SELECT s.id::text,s.challenge_id::text,s.status,tr.passed,tr.failed,tr.total,tr.execution_time_ms,CASE WHEN f.id IS NULL THEN 'PENDING' ELSE 'COMPLETED' END FROM submissions s LEFT JOIN test_results tr ON tr.submission_id=s.id LEFT JOIN followups f ON f.submission_id=s.id WHERE s.id=$1 AND s.candidate_id=$2`, submissionID, candidateID).Scan(&view.SubmissionID, &view.ChallengeID, &view.Status, &passed, &failed, &total, &executionMS, &view.FollowupStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return View{}, ErrNotFound
	}
	if err != nil {
		return View{}, fmt.Errorf("get submission: %w", err)
	}
	if passed != nil {
		view.TestResults = &TestResults{Passed: *passed, Failed: *failed, Total: *total, ExecutionTimeMS: *executionMS}
	}
	return view, nil
}

func (r *PGRepository) SaveFollowup(ctx context.Context, candidateID, submissionID, answer string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var status Status
	var question string
	err = tx.QueryRow(ctx, `SELECT s.status,c.followup_question FROM submissions s JOIN challenges c ON c.id=s.challenge_id WHERE s.id=$1 AND s.candidate_id=$2 FOR UPDATE`, submissionID, candidateID).Scan(&status, &question)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("load follow-up: %w", err)
	}
	if err := CanAnswerFollowup(status, false); err != nil {
		return ErrFollowupUnavailable
	}
	_, err = tx.Exec(ctx, `INSERT INTO followups (submission_id,question,answer) VALUES ($1,$2,$3)`, submissionID, question, answer)
	if uniqueViolation(err) {
		return ErrFollowupUnavailable
	}
	if err != nil {
		return fmt.Errorf("save follow-up: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO assessment_jobs (submission_id) VALUES ($1)`, submissionID); err != nil {
		return fmt.Errorf("queue assessment: %w", err)
	}
	return tx.Commit(ctx)
}
func uniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
