package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"trl-research-backend/internal/entity"
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
		log.Printf("❌ [GetAdminProfile] Missing Authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// Validate and decode JWT
	kp, err := utils.NewEnvKeyProvider()
	if err != nil {
		log.Printf("❌ [GetAdminProfile] Key provider error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key provider error"})
		return
	}
	claims, err := utils.ValidateJWT(tokenString, os.Getenv("JWT_ISSUER"), os.Getenv("JWT_AUDIENCE"), *kp)
	if err != nil {
		log.Printf("❌ [GetAdminProfile] Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Query from Firestore using email from claims (document ID lookup for immediate consistency)
	// Using GetAdminByEmail instead of GetAdminByID to avoid eventual consistency issues
	// with Firestore field queries after updates
	admin, err := h.Repo.GetAdminByEmail(claims.UserEmail)
	if err != nil {
		// Fallback to GetAdminByID if GetAdminByEmail fails
		admin, err = h.Repo.GetAdminByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
			return
		}
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

// 🟢 PATCH /admin/:id
func (h *AdminHandler) UpdateAdminProfileByID(c *gin.Context) {
	// Accept AdminResponse entity type for consistent API format
	var updateReq entity.AdminResponse
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("❌ [UpdateAdminProfile] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")

	// Get existing admin to preserve fields not being updated
	existingAdmin, err := h.Repo.GetAdminByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Update only fields that are provided in the request (non-empty values)
	updateFields := &models.AdminInfo{
		AdminID:               existingAdmin.AdminID,
		AdminPrefix:           existingAdmin.AdminPrefix,
		AdminAcademicPosition: existingAdmin.AdminAcademicPosition,
		AdminFirstName:        existingAdmin.AdminFirstName,
		AdminLastName:         existingAdmin.AdminLastName,
		AdminDepartment:       existingAdmin.AdminDepartment,
		AdminPhoneNumber:      existingAdmin.AdminPhoneNumber,
		AdminEmail:            existingAdmin.AdminEmail,
		AdminPassword:         existingAdmin.AdminPassword,
		CaseID:                existingAdmin.CaseID,
		CreatedAt:             existingAdmin.CreatedAt,
		UpdatedAt:             time.Now(),
	}

	// Only update fields that are provided (non-empty) in the request
	if updateReq.Prefix != "" {
		updateFields.AdminPrefix = updateReq.Prefix
	}
	if updateReq.AcademicPosition != "" {
		updateFields.AdminAcademicPosition = updateReq.AcademicPosition
	}
	if updateReq.FirstName != "" {
		updateFields.AdminFirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		updateFields.AdminLastName = updateReq.LastName
	}
	if updateReq.Department != "" {
		updateFields.AdminDepartment = updateReq.Department
	}
	if updateReq.PhoneNumber != "" {
		updateFields.AdminPhoneNumber = updateReq.PhoneNumber
	}
	
	if err := h.Repo.UpdateAdminByID(id, updateFields); err != nil {
		log.Printf("❌ [UpdateAdminProfile] Error updating admin: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated admin to verify changes
	// Use GetAdminByEmail (document ID lookup) instead of GetAdminByID (field query)
	// to avoid eventual consistency issues with Firestore queries
	updatedAdmin, err := h.Repo.GetAdminByEmail(updateFields.AdminEmail)
	if err != nil {
		// Fallback to GetAdminByID if GetAdminByEmail fails
		updatedAdmin, err = h.Repo.GetAdminByID(id)
		if err != nil {
			log.Printf("⚠️ [UpdateAdminProfile] Failed to verify update: %v", err)
		}
	}

	// Return updated profile in response (using entity.AdminResponse format)
	if updatedAdmin != nil {
		response := updatedAdmin.ToResponse()
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
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
