package handlers

import (
	"log"
	"net/http"

	"trl-research-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type FileDownloadHandler struct {
	GCS *storage.GCSClient
}

func (h *FileDownloadHandler) GetDownloadURL(c *gin.Context) {
	objectPath := c.Query("path")
	log.Printf("[GetDownloadURL] Request for path: %s\n", objectPath)
	if objectPath == "" {
		log.Println("[GetDownloadURL] Missing path parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	url, err := h.GCS.GenerateDownloadSignedURL(objectPath, 15)
	if err != nil {
		log.Printf("[GetDownloadURL] Error generating signed URL: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate signed URL"})
		return
	}

	log.Printf("[GetDownloadURL] Successfully generated download URL for: %s\n", objectPath)
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
