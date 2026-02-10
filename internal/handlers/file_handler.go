package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

type FileHandler struct {
	Repo repository.FileRepository
}

type FileUploadedRequest struct {
	CaseID      string `json:"case_id"`
	Name        string `json:"name"`
	ObjectPath  string `json:"object_path"`
	ContentType string `json:"content_type"`
}

func (h *FileHandler) FileUploaded(c *gin.Context) {

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing"})
		return
	}

	var req FileUploadedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	file := &models.Files{
		ID:          uuid.NewString(),
		CaseID:      req.CaseID,
		Name:        req.Name,
		ObjectPath:  req.ObjectPath,
		Bucket:      "trl-pdf-storage",
		UploadedBy:  userID,
		UploadedAt:  time.Now(),
		ContentType: req.ContentType,
	}

	if err := h.Repo.SaveFile(c, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file)
}
