package handlers

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/storage"
)

// ----------------------------
// Handler
// ----------------------------

type PresignHandler struct {
	GCS *storage.GCSClient
}

// ----------------------------
// Request / Response
// ----------------------------

type PresignUploadRequest struct {
	FileName string `json:"name"`
}

type PresignUploadResponse struct {
	UploadURL  string `json:"upload_url"`
	ObjectPath string `json:"object_path"`
}

// ----------------------------
// POST /trl/presign/upload
// ----------------------------

func (h *PresignHandler) PresignUpload(c *gin.Context) {

	// ----------------------------
	// Parse request
	// ----------------------------
	var req PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.FileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file_name is required",
		})
		return
	}

	// ----------------------------
	// Auth context
	// ----------------------------
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// ----------------------------
	// Sanitize filename
	// ----------------------------
	safeFileName := path.Base(req.FileName)

	// ----------------------------
	// Build TEMP object path
	// attachments/tmp/{userID}/{yyyy-mm-dd}/{filename}
	// ----------------------------
	today := time.Now().UTC().Format("2006-01-02")

	objectPath := fmt.Sprintf(
		"attachments/tmp/%s/%s/%s",
		userID,
		today,
		safeFileName,
	)

	fmt.Println("Object path:", objectPath)

	// ----------------------------
	// Generate signed upload URL
	// ----------------------------
	uploadURL, err := h.GCS.GenerateUploadSignedURL(objectPath, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate upload url",
		})
		return
	}

	fmt.Println("Upload URL:", uploadURL)

	// ----------------------------
	// Response
	// ----------------------------
	c.JSON(http.StatusOK, PresignUploadResponse{
		UploadURL:  uploadURL,
		ObjectPath: objectPath,
	})
}
