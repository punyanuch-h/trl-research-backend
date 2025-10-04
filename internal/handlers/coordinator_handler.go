package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var coordinatorRepo = &repository.CoordinatorRepo{}

func CreateOrUpdateCoordinator(c *gin.Context) {
	var req models.CoordinatorInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := coordinatorRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func GetCoordinator(c *gin.Context) {
	email := c.Param("email")
	res, err := coordinatorRepo.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteCoordinator(c *gin.Context) {
	email := c.Param("email")
	if err := coordinatorRepo.Delete(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": email})
}
