package auth

import (
	"net/http"

	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResetHandler struct {
	AdminRepo repository.AdminRepo
}

type ResetReq struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *ResetHandler) ResetPassword(c *gin.Context) {
	var req ResetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Email == "" || req.OldPassword == "" || len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// verify
	if _, err := h.AdminRepo.Login(req.Email, req.OldPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, string(hash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
