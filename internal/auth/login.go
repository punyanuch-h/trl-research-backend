package auth

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	AdminRepo repository.AdminRepo
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Single call to Login method (which includes password verification)
	admin, err := h.AdminRepo.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	kp, err := utils.NewFileKeyProvider()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "key provider error"})
		return
	}

	expH, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil || expH <= 0 {
		expH = 24
	}

	token, err := utils.GenerateJWT(
		admin.AdminID,
		admin.AdminEmail,
		"admin",
		"", // clientID (fill when available)
		"", // clientName (fill when available)
		os.Getenv("JWT_ISSUER"),
		os.Getenv("JWT_AUDIENCE"),
		"v1",
		time.Duration(expH),
		*kp,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": expH,
	})
}
