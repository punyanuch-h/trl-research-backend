package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

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
	log.Println("claims:", claims)
	// Query from Firestore using user_id from claims
	researcher, err := h.Repo.GetResearcherByID(claims.UserID)
	if err != nil {
		log.Println("Researcher not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Researcher not found"})
		return
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
func (h *ResearcherHandler) UpdateResearcherByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Println("Can not bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateResearcherByID(id, updateData); err != nil {
		log.Println("Update Researcher error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Researcher updated successfully"})
}
