package handlers

import (
	"net/http"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/storage"
)

type PresignHandler struct {
	GCS *storage.GCSClient
}

type PresignRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectPath string `json:"object_path"`
}

// POST /presign/upload
func (h *PresignHandler) PresignUpload(c *gin.Context) {

	var req PresignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if req.ContentType == "" {
		req.ContentType = "application/pdf"
	}

	// generate object path
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID missing or invalid"})
		return
	}
	
	today := time.Now().Format("2006-01-02")
	objectPath := fmt.Sprintf("pdf/%s/%s/%s",  today, userID, req.FileName)

	// generate URL
	url, err := h.GCS.GenerateUploadSignedURL(objectPath, req.ContentType, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := PresignResponse{
		UploadURL: url,
		ObjectPath: objectPath,
	}

	c.JSON(http.StatusOK, res)
}
