package assessment

import (
	"context"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerRepo struct{}

func (handlerRepo) Save(context.Context, string, string, Result, string, string) (StoredResult, error) {
	return StoredResult{}, nil
}
func (handlerRepo) Get(context.Context, string, string) (StoredResult, error) {
	return StoredResult{AssessmentID: "assessment", Result: Result{Score: 82, IntegrityStatus: IntegrityHigh, Feedback: "feedback", Evidence: []Evidence{{Competency: "DEBUGGING", Status: Demonstrated, Explanation: "evidence"}}}}, nil
}
func TestAssessmentHandlerRequiresCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(NewService(handlerRepo{})).Register(router.Group("/api/v1"))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/assessments/assessment", nil))
	if writer.Code != http.StatusForbidden {
		t.Fatalf("got %d", writer.Code)
	}
}
func TestAssessmentHandlerReturnsOwnerView(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set(submission.CandidateContextKey, "candidate"); c.Next() })
	NewHandler(NewService(handlerRepo{})).Register(router.Group("/api/v1"))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/assessments/assessment", nil))
	if writer.Code != http.StatusOK {
		t.Fatalf("got %d: %s", writer.Code, writer.Body.String())
	}
}
