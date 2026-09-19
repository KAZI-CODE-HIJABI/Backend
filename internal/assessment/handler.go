package assessment

import (
	"errors"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/response"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/submission"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler          { return &Handler{service: service} }
func (h *Handler) Register(routes *gin.RouterGroup) { routes.GET("/assessments/:assessment_id", h.get) }
func (h *Handler) get(c *gin.Context) {
	value, ok := c.Get(submission.CandidateContextKey)
	candidateID, isString := value.(string)
	if !ok || !isString || candidateID == "" {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Candidate access is required.")
		return
	}
	stored, err := h.service.Get(c.Request.Context(), candidateID, c.Param("assessment_id"))
	if errors.Is(err, ErrNotFound) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Assessment not found.")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"assessment_id": stored.AssessmentID, "score": stored.Result.Score, "integrity_status": stored.Result.IntegrityStatus, "feedback": stored.Result.Feedback, "evidence": stored.Result.Evidence})
}
