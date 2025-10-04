package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var supporterRepo = &repository.SupporterRepo{}

func CreateOrUpdateSupporter(c *gin.Context) {
	var req models.Supporter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := supporterRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func GetSupporter(c *gin.Context) {
	caseID := c.Param("case_id")
	res, err := supporterRepo.GetByCaseID(caseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteSupporter(c *gin.Context) {
	caseID := c.Param("case_id")
	if err := supporterRepo.Delete(caseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": caseID})
}
