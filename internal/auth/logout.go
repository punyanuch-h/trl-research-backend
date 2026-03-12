package auth

import (
	"net/http"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	RefreshTokenRepo repository.RefreshTokenRepository
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *LogoutHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	tokenHash := utils.HashToken(req.RefreshToken)
	token, err := h.RefreshTokenRepo.GetByHash(tokenHash)
	if err == nil && token != nil {
		_ = h.RefreshTokenRepo.Revoke(token.ID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
