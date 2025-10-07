package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
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