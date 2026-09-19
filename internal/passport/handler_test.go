package passport

import (
	"context"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerRepository struct{}

func (handlerRepository) Get(context.Context, string) (View, error) {
	return View{CandidateID: "candidate", AnonymousID: "KZ-1042", OverallScore: 82}, nil
}
func TestPassportHandlerRequiresCandidate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(NewService(handlerRepository{})).Register(router.Group("/api/v1"))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/passport/me", nil))
	if writer.Code != http.StatusForbidden {
		t.Fatalf("got %d", writer.Code)
	}
}
func TestPassportHandlerReturnsCandidatePassport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set(submission.CandidateContextKey, "candidate"); c.Next() })
	NewHandler(NewService(handlerRepository{})).Register(router.Group("/api/v1"))
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/api/v1/passport/me", nil))
	if writer.Code != http.StatusOK {
		t.Fatalf("got %d: %s", writer.Code, writer.Body.String())
	}
}
