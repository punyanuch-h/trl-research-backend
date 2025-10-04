package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var trlPart2Repo = &repository.AssessmentTrlPart2Repo{}

func CreateOrUpdateTrlPart2(c *gin.Context) {
	var req models.AssessmentTRLPart2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := trlPart2Repo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func GetTrlPart2(c *gin.Context) {
	id := c.Param("id")
	res, err := trlPart2Repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteTrlPart2(c *gin.Context) {
	id := c.Param("id")
	if err := trlPart2Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
