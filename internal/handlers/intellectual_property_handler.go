package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

type IntellectualPropertyHandler struct {
	Repo *repository.IntellectualPropertyRepo
}

// 🟢 GET /ips
func (h *IntellectualPropertyHandler) GetIPAll(c *gin.Context) {
	ips, err := h.Repo.GetIPAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ips)
}

// 🟢 GET /ip/:id
func (h *IntellectualPropertyHandler) GetIPByID(c *gin.Context) {
	id := c.Param("id")
	ip, err := h.Repo.GetIPByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Intellectual Property not found"})
		return
	}
	c.JSON(http.StatusOK, ip)
}

// 🟢 GET /ip/case/:id
func (h *IntellectualPropertyHandler) GetIPByCaseID(c *gin.Context) {
	id := c.Param("id")
	ip, err := h.Repo.GetIPByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Intellectual Property not found"})
		return
	}
	c.JSON(http.StatusOK, ip)
}

// 🟢 POST /ip
func (h *IntellectualPropertyHandler) CreateIP(c *gin.Context) {
	var req models.IntellectualProperty
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateIP(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /ip/:id
func (h *IntellectualPropertyHandler) UpdateIPByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateIPByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Intellectual Property updated successfully"})
}
