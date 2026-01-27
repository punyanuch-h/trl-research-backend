package handlers

import (
	"net/http"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SupportmentHandler struct {
	Repo repository.SupportmentRepository
}

// 🟢 GET /supportments
func (h *SupportmentHandler) GetSupportmentAll(c *gin.Context) {
	supportments, err := h.Repo.GetSupportmentAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, supportments)
}

// 🟢 GET /supportment/:id
func (h *SupportmentHandler) GetSupportmentByID(c *gin.Context) {
	id := c.Param("id")
	supportment, err := h.Repo.GetSupportmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supportment not found"})
		return
	}
	c.JSON(http.StatusOK, supportment)
}

// 🟢 GET /supportment/case/:id
func (h *SupportmentHandler) GetSupportmentByCaseID(c *gin.Context) {
	id := c.Param("id")
	supportment, err := h.Repo.GetSupportmentByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supportment not found"})
		return
	}
	c.JSON(http.StatusOK, supportment)
}

// 🟢 POST /supportment
func (h *SupportmentHandler) CreateSupportment(c *gin.Context) {
	var req models.Supportments
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateSupportment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /supportment/:id
func (h *SupportmentHandler) UpdateSupportmentByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateSupportmentByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supportment updated successfully"})
}
