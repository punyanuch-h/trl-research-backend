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
    Repo *repository.FileRepo
}

type FileUploadedRequest struct {
    FileName        string `json:"file_name"`
    ObjectPath      string `json:"object_path"`
    ContentType     string `json:"content_type"`
    BelongsToCaseID string `json:"belongs_to_case_id"`
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

    file := &models.FileMetadata{
        ID:              uuid.NewString(),
        FileName:        req.FileName,
        ObjectPath:      req.ObjectPath,
        Bucket:          "trl-pdf-storage",
        UploadedBy:      userID,
        UploadedAt:      time.Now(),
        ContentType:     req.ContentType,
        BelongsToCaseID: req.BelongsToCaseID,
    }

    if err := h.Repo.SaveFile(c, file); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, file)
}
