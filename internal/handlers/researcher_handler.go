package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var researcherRepo = &repository.ResearcherRepo{}

func CreateOrUpdateResearcher(c *gin.Context) {
	var req models.ResearcherInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CreatedAt = time.Now()

	if err := researcherRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func GetResearcher(c *gin.Context) {
	id := c.Param("id")
	res, err := researcherRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteResearcher(c *gin.Context) {
	id := c.Param("id")
	if err := researcherRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
