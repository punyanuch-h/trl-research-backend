package handlers

import (
	"net/http"
	"os"
	"strings"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Repo *repository.AdminRepo
}

// 🟢 GET /admins
func (h *AdminHandler) GetAllAdmins(c *gin.Context) {
	admins, err := h.Repo.GetAdminAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, admins)
}

// 🟢 GET /admin/:id
func (h *AdminHandler) GetAdminByID(c *gin.Context) {
	id := c.Param("id")
	admin, err := h.Repo.GetAdminByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}
	c.JSON(http.StatusOK, admin)
}

// 🟢 GET /admin/profile
func (h *AdminHandler) GetAdminProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// Validate and decode JWT
	kp, err := utils.NewEnvKeyProvider()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key provider error"})
		return
	}
	claims, err := utils.ValidateJWT(tokenString, os.Getenv("JWT_ISSUER"), os.Getenv("JWT_AUDIENCE"), *kp)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	// Query from Firestore using user_id from claims
	admin, err := h.Repo.GetAdminByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Convert to response DTO
	response := admin.ToResponse()
	c.JSON(http.StatusOK, response)
}

// 🟢 POST /admin
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req models.AdminInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAdmin(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// update password
func (h *AdminHandler) UpdatePassword(c *gin.Context) {
	var req models.AdminInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// Login
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.AdminInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// delete admin
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	email := c.Param("email")
	if err := h.Repo.DeleteAdmin(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}