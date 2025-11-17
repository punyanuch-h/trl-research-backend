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

type ResearcherHandler struct {
	Repo *repository.ResearcherRepo
}

// 🟢 GET /researchers
func (h *ResearcherHandler) GetResearcherAll(c *gin.Context) {
	researchers, err := h.Repo.GetResearcherAll()
	if err != nil {
		log.Println("Get Researcher All error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, researchers)
}

// 🟢 GET /researcher/:id
func (h *ResearcherHandler) GetResearcherByID(c *gin.Context) {
	id := c.Param("id")
	researcher, err := h.Repo.GetResearcherByID(id)
	if err != nil {
		log.Println("Researcher not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
		return
	}
	c.JSON(http.StatusOK, researcher)
}

// 🟢 GET /researcher/case/:id
func (h *ResearcherHandler) GetResearcherByCaseID(c *gin.Context) {
	id := c.Param("id")
	researcher, err := h.Repo.GetResearcherByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
		return
	}
	c.JSON(http.StatusOK, researcher)
}

// 🟢 GET /researcher/profile
func (h *ResearcherHandler) GetResearcherProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Printf("❌ [GetResearcherProfile] Missing Authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// Validate and decode JWT
	kp, err := utils.NewEnvKeyProvider()
	if err != nil {
		log.Printf("❌ [GetResearcherProfile] Key provider error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Key provider error"})
		return
	}
	claims, err := utils.ValidateJWT(tokenString, os.Getenv("JWT_ISSUER"), os.Getenv("JWT_AUDIENCE"), *kp)
	if err != nil {
		log.Printf("❌ [GetResearcherProfile] Invalid token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Query from Firestore using user_id from claims (document ID lookup for immediate consistency)
	// For researchers, document ID = researcher_id, so we can use GetResearcherByIDDirect
	// which uses document ID lookup instead of field query
	researcher, err := h.Repo.GetResearcherByIDDirect(claims.UserID)
	if err != nil {
		// Fallback to field query if document ID lookup fails
		researcher, err = h.Repo.GetResearcherByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
			return
		}
	}

	// Convert to response DTO
	response := researcher.ToResponse()
	c.JSON(http.StatusOK, response)
}

// 🟢 POST /researcher
func (h *ResearcherHandler) CreateResearcher(c *gin.Context) {
	var req models.ResearcherInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Can not bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateResearcher(&req); err != nil {
		log.Println("Create Researcher error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /researcher/:id
func (h *ResearcherHandler) UpdateResearcherProfileByID(c *gin.Context) {
	// Accept ResearcherResponse entity type for consistent API format
	var updateReq entity.ResearcherResponse
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		log.Printf("❌ [UpdateResearcherProfile] Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")

	// Get existing researcher to preserve fields not being updated
	existingResearcher, err := h.Repo.GetResearcherByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
		return
	}

	// Update only fields that are provided in the request (non-empty values)
	updateFields := &models.ResearcherInfo{
		ResearcherID:               existingResearcher.ResearcherID,
		AdminID:                    existingResearcher.AdminID,
		ResearcherPrefix:           existingResearcher.ResearcherPrefix,
		ResearcherAcademicPosition: existingResearcher.ResearcherAcademicPosition,
		ResearcherFirstName:        existingResearcher.ResearcherFirstName,
		ResearcherLastName:         existingResearcher.ResearcherLastName,
		ResearcherDepartment:       existingResearcher.ResearcherDepartment,
		ResearcherPhoneNumber:      existingResearcher.ResearcherPhoneNumber,
		ResearcherEmail:            existingResearcher.ResearcherEmail,
		ResearcherPassword:         existingResearcher.ResearcherPassword,
		CreatedAt:                  existingResearcher.CreatedAt,
		UpdatedAt:                  time.Now(),
	}

	// Only update fields that are provided (non-empty) in the request
	if updateReq.Prefix != "" {
		updateFields.ResearcherPrefix = updateReq.Prefix
	}
	if updateReq.AcademicPosition != "" {
		updateFields.ResearcherAcademicPosition = updateReq.AcademicPosition
	}
	if updateReq.FirstName != "" {
		updateFields.ResearcherFirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		updateFields.ResearcherLastName = updateReq.LastName
	}
	if updateReq.Department != "" {
		updateFields.ResearcherDepartment = updateReq.Department
	}
	if updateReq.PhoneNumber != "" {
		updateFields.ResearcherPhoneNumber = updateReq.PhoneNumber
	}
	
	if err := h.Repo.UpdateResearcherByID(id, updateFields); err != nil {
		log.Printf("❌ [UpdateResearcherProfile] Error updating researcher: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated researcher to verify changes
	// Use GetResearcherByIDDirect (document ID lookup) instead of GetResearcherByID (field query)
	// to avoid eventual consistency issues with Firestore queries
	updatedResearcher, err := h.Repo.GetResearcherByIDDirect(id)
	if err != nil {
		// Fallback to GetResearcherByID if GetResearcherByIDDirect fails
		updatedResearcher, err = h.Repo.GetResearcherByID(id)
		if err != nil {
			log.Printf("⚠️ [UpdateResearcherProfile] Failed to verify update: %v", err)
		}
	}

	// Return updated profile in response (using entity.ResearcherResponse format)
	if updatedResearcher != nil {
		response := updatedResearcher.ToResponse()
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Researcher updated successfully"})
	}
}
