package response

import "github.com/gin-gonic/gin"

// Error is the proposed shared envelope for infrastructure errors.
// Domain error bodies must be agreed in the draft API contract before use.
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}
