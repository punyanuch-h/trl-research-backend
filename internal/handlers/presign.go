package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"trl-research-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type PresignHandler struct {
	GCS *storage.GCSClient
}

type PresignUploadRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type PresignUploadResponse struct {
	UploadURL  string `json:"upload_url"`
	ObjectPath string `json:"object_path"`
}

func (h *PresignHandler) PresignUpload(c *gin.Context) {
	log.Println("[PresignUpload] Incoming request")
	var req PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[PresignUpload] JSON bind error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.FileName == "" {
		log.Println("[PresignUpload] Missing file_name in request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name is required",
		})
		return
	}

	// Default to application/octet-stream as requested by user
	if req.ContentType == "" {
		req.ContentType = "application/octet-stream"
	}

	// Auth context
	userID := c.GetString("userID")
	if userID == "" {
		log.Println("[PresignUpload] Unauthorized: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	log.Printf("[PresignUpload] Request for file: %s (%s) by user: %s\n", req.FileName, req.ContentType, userID)

	// Sanitize filename
	safeFileName := path.Base(req.FileName)
	if safeFileName == "" || safeFileName == "." {
		log.Printf("[PresignUpload] Invalid filename after sanitization: %s\n", req.FileName)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid filename",
		})
		return
	}

	// Build TEMP object path
	// attachments/tmp/{userID}/{yyyy-mm-dd}/{filename}
	today := time.Now().UTC().Format("2006-01-02")

	objectPath := fmt.Sprintf(
		"attachments/tmp/%s/%s/%s",
		userID,
		today,
		safeFileName,
	)

	log.Printf("[PresignUpload] Generating signed URL for path: %s with Content-Type: %s\n", objectPath, req.ContentType)

	uploadURL, err := h.GCS.GenerateUploadSignedURL(objectPath, req.ContentType, 15)
	if err != nil {
		log.Printf("[PresignUpload] Error generating signed URL: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate upload url",
		})
		return
	}

	log.Printf("[PresignUpload] Successfully generated signed URL for: %s\n", objectPath)

	c.JSON(http.StatusOK, PresignUploadResponse{
		UploadURL:  uploadURL,
		ObjectPath: objectPath,
	})
}
