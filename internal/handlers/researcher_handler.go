package handlers

import (
	"log"
	"net/http"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

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
