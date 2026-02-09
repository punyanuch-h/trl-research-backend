package handlers

import (
	"net/http"

	"trl-research-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type FileDownloadHandler struct {
	GCS *storage.GCSClient
}

func (h *FileDownloadHandler) GetDownloadURL(c *gin.Context) {
	objectPath := c.Query("path")
	if objectPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	url, err := h.GCS.GenerateDownloadSignedURL(objectPath, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
