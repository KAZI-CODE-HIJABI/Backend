package submission

import (
	"errors"
	"net/http"

	"github.com/KAZI-CODE-HIJABI/Backend/internal/response"
	"github.com/gin-gonic/gin"
)

// CandidateContextKey is the contract between the auth middleware and this module.
// Auth must set it only after validating a CANDIDATE access token.
const CandidateContextKey = "candidate_id"

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) Register(routes *gin.RouterGroup) {
	routes.POST("/submissions", h.create)
	routes.GET("/submissions/:submission_id", h.get)
	routes.POST("/submissions/:submission_id/followup", h.followup)
}
func (h *Handler) candidateID(c *gin.Context) (string, bool) {
	id, ok := c.Get(CandidateContextKey)
	value, isString := id.(string)
	return value, ok && isString && value != ""
}
func (h *Handler) create(c *gin.Context) {
	candidateID, ok := h.candidateID(c)
	if !ok {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Candidate access is required.")
		return
	}
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "challenge_id, code, and explanation are required.")
		return
	}
	if err := input.Validate(); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	view, err := h.service.Create(c.Request.Context(), candidateID, input)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"submission_id": view.SubmissionID, "status": view.Status})
}
func (h *Handler) get(c *gin.Context) {
	candidateID, ok := h.candidateID(c)
	if !ok {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Candidate access is required.")
		return
	}
	view, err := h.service.Get(c.Request.Context(), candidateID, c.Param("submission_id"))
	if err != nil {
		h.writeError(c, err)
		return
	}
	body := gin.H{"submission_id": view.SubmissionID, "challenge_id": view.ChallengeID, "status": view.Status, "followup_status": view.FollowupStatus}
	if view.TestResults != nil {
		body["test_results"] = gin.H{"passed": view.TestResults.Passed, "failed": view.TestResults.Failed, "total": view.TestResults.Total, "execution_time": float64(view.TestResults.ExecutionTimeMS) / 1000}
	}
	c.JSON(http.StatusOK, body)
}
func (h *Handler) followup(c *gin.Context) {
	candidateID, ok := h.candidateID(c)
	if !ok {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Candidate access is required.")
		return
	}
	var request struct {
		Answer string `json:"answer"`
	}
	if err := c.ShouldBindJSON(&request); err != nil || request.Answer == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "answer is required.")
		return
	}
	if err := h.service.AnswerFollowup(c.Request.Context(), candidateID, c.Param("submission_id"), request.Answer); err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"submission_id": c.Param("submission_id"), "followup_status": "COMPLETED"})
}
func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(c, 404, "NOT_FOUND", "Submission not found.")
	case errors.Is(err, ErrActiveSubmission), errors.Is(err, ErrFollowupUnavailable):
		response.Error(c, 409, "CONFLICT", err.Error())
	default:
		response.Error(c, 500, "INTERNAL_ERROR", "An unexpected error occurred.")
	}
}
