package handlers

import (
	"net/http"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

var adminRepo = &repository.AdminRepo{}

func CreateAdmin(c *gin.Context) {
	var req models.AdminInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if req.AdminID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin_id is required"})
		return
	}
	req.CreatedAt = time.Now()

	if err := adminRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, req)
	c.Set("log", "Admin created successfully")
}

func GetAdmin(c *gin.Context) {
	id := c.Param("id")
	res, err := adminRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
	c.Set("log", "Admin fetched successfully")
}

func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	if err := adminRepo.DeleteAdmin(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
	c.Set("log", "Admin deleted successfully")
}
