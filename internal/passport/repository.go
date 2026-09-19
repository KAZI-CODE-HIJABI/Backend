package passport

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepository struct{ pool *pgxpool.Pool }

func NewPGRepository(pool *pgxpool.Pool) *PGRepository { return &PGRepository{pool: pool} }
func (r *PGRepository) Get(ctx context.Context, candidateID string) (View, error) {
	var view View
	var passed, failed, total int
	err := r.pool.QueryRow(ctx, `SELECT cp.candidate_id::text,cp.anonymous_id,a.score,COALESCE(tr.passed,0),COALESCE(tr.failed,0),COALESCE(tr.total,0) FROM candidate_profiles cp JOIN submissions s ON s.candidate_id=cp.candidate_id JOIN assessments a ON a.submission_id=s.id LEFT JOIN test_results tr ON tr.submission_id=s.id WHERE cp.candidate_id=$1 ORDER BY a.created_at DESC LIMIT 1`, candidateID).Scan(&view.CandidateID, &view.AnonymousID, &view.OverallScore, &passed, &failed, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return View{}, ErrNotFound
	}
	if err != nil {
		return View{}, fmt.Errorf("get passport: %w", err)
	}
	view.Category = "SOFTWARE_ENGINEERING"
	view.EvidenceSummary = Summary{ChallengesCompleted: 1, TestsPassed: fmt.Sprintf("%d/%d", passed, total), ExplanationSubmitted: true, FollowupCompleted: true}
	rows, err := r.pool.Query(ctx, `SELECT sk.name,cs.status FROM candidate_skills cs JOIN skills sk ON sk.id=cs.skill_id WHERE cs.candidate_id=$1 ORDER BY sk.name`, candidateID)
	if err != nil {
		return View{}, fmt.Errorf("get passport skills: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.Name, &skill.Status); err != nil {
			return View{}, err
		}
		view.Skills = append(view.Skills, skill)
	}
	if err := rows.Err(); err != nil {
		return View{}, err
	}
	return view, nil
}
