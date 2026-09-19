package submission

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerRepository struct{}

func (handlerRepository) Create(context.Context, string, CreateInput) (View, error) {
	return View{SubmissionID: "submission", Status: Queued, FollowupStatus: "PENDING"}, nil
}
func (handlerRepository) Get(context.Context, string, string) (View, error) {
	return View{SubmissionID: "submission", ChallengeID: "challenge", Status: Completed, FollowupStatus: "PENDING", TestResults: &TestResults{Passed: 8, Failed: 2, Total: 10, ExecutionTimeMS: 1820}}, nil
}
func (handlerRepository) SaveFollowup(context.Context, string, string, string) error { return nil }
func candidateMiddleware(c *gin.Context)                                             { c.Set(CandidateContextKey, "candidate"); c.Next() }
func TestCreateHandlerMatchesContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(candidateMiddleware)
	NewHandler(NewService(handlerRepository{})).Register(router.Group("/api/v1"))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", strings.NewReader(`{"challenge_id":"challenge","code":"code","explanation":"reason"}`))
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	if writer.Code != http.StatusAccepted {
		t.Fatalf("got %d: %s", writer.Code, writer.Body.String())
	}
}
func TestCreateHandlerRejectsMissingCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(NewService(handlerRepository{})).Register(router.Group("/api/v1"))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest(http.MethodPost, "/api/v1/submissions", nil))
	if writer.Code != http.StatusForbidden {
		t.Fatalf("got %d", writer.Code)
	}
}
