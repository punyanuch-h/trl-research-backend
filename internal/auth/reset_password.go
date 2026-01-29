package auth

import (
	"net/http"

	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type ResetHandler struct {
	AdminRepo repository.AdminRepository
	ResearcherRepo repository.ResearcherRepository
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
	if req.Email == "" || req.OldPassword == "" || len(req.NewPassword) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Verify and update for admin
    if _, err := h.AdminRepo.Login(req.Email, req.OldPassword); err == nil {
        if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, req.NewPassword); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "password updated"})
        return
    }

    // Verify and update for researcher
    if _, err := h.ResearcherRepo.Login(req.Email, req.OldPassword); err == nil {
        if err := h.ResearcherRepo.UpdatePasswordByEmail(req.Email, req.NewPassword); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "password updated"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
