package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

type SupporterHandler struct {
	Repo *repository.SupporterRepo
}

// 🟢 GET /supporters
func (h *SupporterHandler) GetSupporterAll(c *gin.Context) {
	supporters, err := h.Repo.GetSupporterAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, supporters)
}

// 🟢 GET /supporter/:id
func (h *SupporterHandler) GetSupporterByID(c *gin.Context) {
	id := c.Param("id")
	supporter, err := h.Repo.GetSupporterByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supporter not found"})
		return
	}
	c.JSON(http.StatusOK, supporter)
}

// 🟢 GET /supporter/case/:id
func (h *SupporterHandler) GetSupporterByCaseID(c *gin.Context) {
	id := c.Param("id")
	supporter, err := h.Repo.GetSupporterByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supporter not found"})
		return
	}
	c.JSON(http.StatusOK, supporter)
}

// 🟢 POST /supporter
func (h *SupporterHandler) CreateSupporter(c *gin.Context) {
	var req models.Supporter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateSupporter(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /supporter/:id
func (h *SupporterHandler) UpdateSupporterByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateSupporterByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supporter updated successfully"})
}
