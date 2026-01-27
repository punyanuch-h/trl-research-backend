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
	Repo repository.ResearcherRepository
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

	// Query from Postgres using user_id from claims
	researcher, err := h.Repo.GetResearcherByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
		return
	}

	// Convert to response DTO
	response := researcher.ToResponse()
	c.JSON(http.StatusOK, response)
}

// 🟢 POST /researcher
func (h *ResearcherHandler) CreateResearcher(c *gin.Context) {
	var req models.Researchers
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
	updateFields := &models.Researchers{
		ID:               existingResearcher.ID,
		Prefix:           existingResearcher.Prefix,
		AcademicPosition: existingResearcher.AcademicPosition,
		FirstName:        existingResearcher.FirstName,
		LastName:         existingResearcher.LastName,
		Department:       existingResearcher.Department,
		PhoneNumber:      existingResearcher.PhoneNumber,
		Email:            existingResearcher.Email,
		Password:         existingResearcher.Password,
		CreatedAt:        existingResearcher.CreatedAt,
		UpdatedAt:        time.Now(),
	}

	// Only update fields that are provided (non-empty) in the request
	if updateReq.Prefix != "" {
		updateFields.Prefix = updateReq.Prefix
	}
	if updateReq.AcademicPosition != "" {
		updateFields.AcademicPosition = updateReq.AcademicPosition
	}
	if updateReq.FirstName != "" {
		updateFields.FirstName = updateReq.FirstName
	}
	if updateReq.LastName != "" {
		updateFields.LastName = updateReq.LastName
	}
	if updateReq.Department != "" {
		updateFields.Department = updateReq.Department
	}
	if updateReq.PhoneNumber != "" {
		updateFields.PhoneNumber = updateReq.PhoneNumber
	}

	if err := h.Repo.UpdateResearcherByID(id, updateFields); err != nil {
		log.Printf("❌ [UpdateResearcherProfile] Error updating researcher: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated researcher to verify changes
	updatedResearcher, err := h.Repo.GetResearcherByID(id)
	if err != nil {
		log.Printf("⚠️ [UpdateResearcherProfile] Failed to verify update: %v", err)
	}

	// Return updated profile in response (using entity.ResearcherResponse format)
	if updatedResearcher != nil {
		response := updatedResearcher.ToResponse()
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Researcher updated successfully"})
	}
}
