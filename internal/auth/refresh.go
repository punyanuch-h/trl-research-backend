package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type RefreshHandler struct {
	RefreshTokenRepo repository.RefreshTokenRepository
	KeyProvider      utils.IKeyProvider
	Cfg              config.Config
	AdminRepo        repository.AdminRepository
	ResearcherRepo   repository.ResearcherRepository
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *RefreshHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	// 1️⃣ Hash and find the token in DB
	tokenHash := utils.HashToken(req.RefreshToken)
	storedToken, err := h.RefreshTokenRepo.GetByHash(tokenHash)
	if err != nil || storedToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// 2️⃣ Rotation & Theft Detection
	if storedToken.RevokedAt != nil {
		// If this token was already replaced, someone might be reuse an old token (theft attempt)
		if storedToken.ReplacedByToken != "" {
			fmt.Printf("⚠️ Potential token theft detected for user %s. Revoking all sessions.\n", storedToken.UserID)
			_ = h.RefreshTokenRepo.RevokeAllForUser(storedToken.UserID)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
		return
	}

	if storedToken.ExpiryAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	// 3️⃣ Identify the user to get current email and temp status
	var userEmail string
	var isTemp bool
	if storedToken.UserType == "admin" {
		user, err := h.AdminRepo.GetAdminByID(storedToken.UserID)
		if err != nil || user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		userEmail = user.Email
		isTemp = user.PasswordIsTemp
	} else {
		user, err := h.ResearcherRepo.GetResearcherByID(storedToken.UserID)
		if err != nil || user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		userEmail = user.Email
		isTemp = user.PasswordIsTemp
	}

	// 4️⃣ Generate New Access Token
	kp := h.KeyProvider
	if kp == nil {
		kp, _ = utils.NewEnvKeyProvider()
	}

	ttl := h.Cfg.GetJWTExpiry()
	if isTemp {
		ttl = h.Cfg.GetJWTExpiryTemp()
	}

	newToken, err := utils.GenerateJWT(
		storedToken.UserID,
		userEmail,
		storedToken.UserType,
		"", "",
		os.Getenv("JWT_ISSUER"),
		os.Getenv("JWT_AUDIENCE"),
		"v1",
		isTemp,
		ttl,
		kp,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// 5️⃣ Rotate Refresh Token
	newRefreshTokenStr, err := utils.GenerateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	newRefreshExpiry := h.Cfg.GetRefreshTokenExpiry()
	newRefreshTokenModel := &models.RefreshToken{
		UserID:    storedToken.UserID,
		TokenHash: utils.HashToken(newRefreshTokenStr),
		ExpiryAt:  time.Now().Add(newRefreshExpiry),
		UserType:  storedToken.UserType,
	}

	if err := h.RefreshTokenRepo.Create(newRefreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save refresh token"})
		return
	}

	// Invalidate old token
	now := time.Now()
	storedToken.RevokedAt = &now
	storedToken.ReplacedByToken = newRefreshTokenModel.ID
	if err := h.RefreshTokenRepo.Update(storedToken); err != nil {
		fmt.Printf("⚠️ Error updating old refresh token: %v\n", err)
	}

	// 6️⃣ Respond
	c.JSON(http.StatusOK, gin.H{
		"token":         newToken,
		"refresh_token": newRefreshTokenStr,
		"expires_in":    int(ttl.Minutes()),
		"unit":          "minutes",
	})
}
