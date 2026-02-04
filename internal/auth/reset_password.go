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
	// 1. Get authenticated user info from context
	userEmail := c.GetString("userEmail")
	userRole := c.GetString("role")

	if userEmail == "" || userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ResetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Validate email if provided in request matches the authenticated user
	if req.Email != "" && req.Email != userEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "email does not match authenticated user"})
		return
	}

	// Use authenticated email for consistency
	emailToUse := userEmail

	if req.OldPassword == "" || len(req.NewPassword) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body or password too short"})
		return
	}

	// 3. Reset password based on role
	if userRole == "admin" {
		if _, err := h.AdminRepo.Login(emailToUse, req.OldPassword); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
			return
		}
		if err := h.AdminRepo.UpdatePasswordByEmail(emailToUse, req.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
		return
	}

	if userRole == "researcher" {
		if _, err := h.ResearcherRepo.Login(emailToUse, req.OldPassword); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
			return
		}
		if err := h.ResearcherRepo.UpdatePasswordByEmail(emailToUse, req.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role or unauthorized access"})
}

