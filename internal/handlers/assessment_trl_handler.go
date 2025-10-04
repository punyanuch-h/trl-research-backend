package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var trlRepo = &repository.AssessmentTrlRepo{}

func CreateOrUpdateTrl(c *gin.Context) {
	var req models.AssessmentTRL
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := trlRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func GetTrl(c *gin.Context) {
	caseID := c.Param("case_id")
	res, err := trlRepo.GetByCaseID(caseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteTrl(c *gin.Context) {
	caseID := c.Param("case_id")
	if err := trlRepo.Delete(caseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": caseID})
}
