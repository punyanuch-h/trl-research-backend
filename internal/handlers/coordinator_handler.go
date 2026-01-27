package handlers

import (
	"net/http"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type CoordinatorHandler struct {
	Repo repository.CoordinatorRepository
}

// 🟢 GET /coordinators
func (h *CoordinatorHandler) GetCoordinatorAll(c *gin.Context) {
	coordinators, err := h.Repo.GetCoordinatorAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coordinators)
}

// 🟢 GET /coordinator/:email
func (h *CoordinatorHandler) GetCoordinatorByEmail(c *gin.Context) {
	email := c.Param("email")
	coordinator, err := h.Repo.GetCoordinatorByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coordinator not found"})
		return
	}
	c.JSON(http.StatusOK, coordinator)
}

// 🟢 GET /coordinator/case/:id
func (h *CoordinatorHandler) GetCoordinatorByCaseID(c *gin.Context) {
	id := c.Param("id")
	coordinator, err := h.Repo.GetCoordinatorByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coordinator not found"})
		return
	}
	c.JSON(http.StatusOK, coordinator)
}

// 🟢 POST /coordinator
func (h *CoordinatorHandler) CreateCoordinator(c *gin.Context) {
	var req models.Coordinators
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateCoordinator(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /coordinator/:email
func (h *CoordinatorHandler) UpdateCoordinatorByEmail(c *gin.Context) {
	email := c.Param("email")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateCoordinatorByEmail(email, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coordinator updated successfully"})
}
