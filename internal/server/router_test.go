package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeDB struct{ err error }

func (f fakeDB) Ping(context.Context) error { return f.err }

func TestHealthAndReadiness(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, path string
		dbErr      error
		code       int
	}{
		{"alive without database", "/healthz", errors.New("secret-dsn"), 200},
		{"ready", "/readyz", nil, 200},
		{"database down", "/readyz", errors.New("secret-dsn"), 503},
		{"unimplemented endpoint", "/api/v1/candidates/me", nil, 404},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			Router(fakeDB{tc.dbErr}, zap.NewNop()).ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
			if w.Code != tc.code {
				t.Fatalf("got %d: %s", w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "secret-dsn") {
				t.Fatal("database error leaked")
			}
		})
	}
}
