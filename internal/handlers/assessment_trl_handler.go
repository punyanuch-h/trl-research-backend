package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

type AssessmentTrlHandler struct {
	Repo *repository.AssessmentTrlRepo
}

// 🟢 GET /assessments
func (h *AssessmentTrlHandler) GetAssessmentTrlAll(c *gin.Context) {
	assessments, err := h.Repo.GetAssessmentTrlAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessments)
}

// 🟢 GET /assessment/:id
func (h *AssessmentTrlHandler) GetAssessmentTrlByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.Repo.GetAssessmentTrlByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment TRL not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// 🟢 POST /assessment
func (h *AssessmentTrlHandler) CreateAssessmentTrl(c *gin.Context) {
	var req models.AssessmentTrl
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAssessmentTrl(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /assessment/:id
func (h *AssessmentTrlHandler) UpdateAssessmentTrlByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateAssessmentTrlByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment TRL updated successfully"})
}
