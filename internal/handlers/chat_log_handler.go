package handlers

import (
	"net/http"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ChatLogHandler struct {
	Repo repository.ChatLogRepository
}

type CreateChatLogRequest struct {
	AdminID      *string        `json:"admin_id,omitempty"`
	ResearcherID *string        `json:"researcher_id,omitempty"`
	History      datatypes.JSON `json:"history" binding:"required"`
}

func (h *ChatLogHandler) CreateChatLog(c *gin.Context) {
	var req CreateChatLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatLog := models.ChatLogs{
		AdminID:      req.AdminID,
		ResearcherID: req.ResearcherID,
		History:      req.History,
	}

	if err := h.Repo.CreateChatLog(&chatLog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatLog)
}
