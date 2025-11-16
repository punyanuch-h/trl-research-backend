package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"
)

type FileDownloadHandler struct {
	FileRepo *repository.FileRepo
	GCS      *storage.GCSClient
}

func (h *FileDownloadHandler) GetDownloadURL(c *gin.Context) {
	fileID := c.Param("fileID")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileID is required"})
		return
	}

	userID := c.GetString("userID")
	role := c.GetString("role")

	// Load file metadata
	file, err := h.FileRepo.GetFileByID(c, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// -----------------------------
	// ⭐ Permission Checking
	// -----------------------------
	if role != "admin" && file.UploadedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not allowed to download this file"})
		return
	}

	// Generate signed URL
	url, err := h.GCS.GenerateDownloadSignedURL(file.ObjectPath, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
