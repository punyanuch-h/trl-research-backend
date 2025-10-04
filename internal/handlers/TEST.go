package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

var trlQuestionRepo = &repository.TrlQuestionRepo{}

// CreateOrUpdate → POST /trl_questions
func CreateOrUpdateTrlQuestion(c *gin.Context) {
	var req models.TrlQuestion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := trlQuestionRepo.CreateOrUpdate(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// GetOne → GET /trl_questions/:assessment_id/:question_code
func GetTrlQuestion(c *gin.Context) {
	assessmentID := c.Param("assessment_id")
	questionCode := c.Param("question_code")

	res, err := trlQuestionRepo.Get(assessmentID, questionCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetByAssessment → GET /trl_questions/:assessment_id
func GetTrlQuestionsByAssessment(c *gin.Context) {
	assessmentID := c.Param("assessment_id")

	res, err := trlQuestionRepo.GetByAssessment(assessmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Delete → DELETE /trl_questions/:assessment_id/:question_code
func DeleteTrlQuestion(c *gin.Context) {
	assessmentID := c.Param("assessment_id")
	questionCode := c.Param("question_code")

	if err := trlQuestionRepo.Delete(assessmentID, questionCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": assessmentID + "_" + questionCode})
}
