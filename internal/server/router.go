package server

import (
	"context"
	"net/http"
	"time"

	"github.com/KAZI-CODE-HIJABI/Backend/internal/middleware"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DatabasePinger interface{ Ping(context.Context) error }

func Router(db DatabasePinger, log *zap.Logger) *gin.Engine {
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.HandleMethodNotAllowed = true
	r.Use(middleware.Logging(log))
	r.Use(func(c *gin.Context) {
		defer func() {
			if recover() != nil {
				log.Error("request_panic")
				c.Abort()
				response.Error(c, 500, "INTERNAL_ERROR", "An unexpected error occurred.")
			}
		}()
		c.Next()
	})
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			response.Error(c, 503, "NOT_READY", "Database is unavailable.")
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})
	r.NoRoute(func(c *gin.Context) { response.Error(c, 404, "NOT_FOUND", "Route not found.") })
	r.NoMethod(func(c *gin.Context) { response.Error(c, 405, "METHOD_NOT_ALLOWED", "Method not allowed.") })
	// Register implemented business modules under r.Group("/api/v1").
	return r
}
