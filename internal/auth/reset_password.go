package auth

import (
	"net/http"

	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type ResetHandler struct {
	AdminRepo      repository.AdminRepository
	ResearcherRepo repository.ResearcherRepository
}

type ResetReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *ResetHandler) ResetPassword(c *gin.Context) {
	// 1. Get authenticated user info from context claims
	claims, err := GetMiddleware(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: invalid claims"})
		return
	}
	userEmail := claims.UserEmail
	userRole := claims.Role
	isTemp := claims.IsTemp

	if userEmail == "" || userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req ResetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Validate new password length (matching frontend)
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	// 3. If user is NOT temporary, the old password must be provided and correct.
	// For temporary users, the valid JWT is sufficient authorization.
	if !isTemp {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "old password is required"})
			return
		}
		// Verify old password for the correct role
		if userRole == "admin" {
			if _, err := h.AdminRepo.Login(userEmail, req.OldPassword); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
				return
			}
		} else if userRole == "researcher" {
			if _, err := h.ResearcherRepo.Login(userEmail, req.OldPassword); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
				return
			}
		}
	}

	// 4. Update the password and set PasswordIsTemp to false
	if userRole == "admin" {
		if err := h.AdminRepo.UpdatePasswordByEmail(userEmail, req.NewPassword, false, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
		return
	}

	if userRole == "researcher" {
		if err := h.ResearcherRepo.UpdatePasswordByEmail(userEmail, req.NewPassword, false, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "password updated"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role or unauthorized access"})
}
