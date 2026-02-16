package handlers

import (
	"net/http"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type ChatLogHandler struct {
	Repo repository.ChatLogRepository
}

func (h *ChatLogHandler) CreateChatLog(c *gin.Context) {
	var req models.ChatLogs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateChatLog(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}
