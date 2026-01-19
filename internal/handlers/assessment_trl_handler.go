package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type AssessmentTrlHandler struct {
	Repo repository.AssessmentTrlRepository
	GCS  *storage.GCSClient
}

// 🟢 GET /assessments
func (h *AssessmentTrlHandler) GetAssessmentTrlAll(c *gin.Context) {
	assessments, err := h.Repo.GetAssessmentTrlAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessments)
}

// 🟢 GET /assessment/:id
func (h *AssessmentTrlHandler) GetAssessmentTrlByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.Repo.GetAssessmentTrlByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment TRL not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// 🟢 GET /assessment/case/:id
func (h *AssessmentTrlHandler) GetAssessmentTrlByCaseID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.Repo.GetAssessmentTrlByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment TRL not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// 🟢 POST /assessment
func (h *AssessmentTrlHandler) CreateAssessmentTrl(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		var req models.AssessmentTrl
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		form, _ := c.MultipartForm()
		today := time.Now().Format("2006-01-02")
		userID := c.GetString("userID")
		if userID == "" {
			userID = "unknown_user"
		}

		// Helper to upload files for a specific key and return as datatypes.JSON
		uploadFiles := func(key string) datatypes.JSON {
			files := form.File[key]
			var paths []string
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					continue
				}
				defer file.Close()

				objectPath := fmt.Sprintf("assessment_attachments/%s/%s/%s/%s", today, req.CaseID, key, fileHeader.Filename)
				if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err == nil {
					paths = append(paths, objectPath)
				}
			}
			jsonData, _ := json.Marshal(paths)
			return datatypes.JSON(jsonData)
		}

		// Process all attachments
		req.Rq1Attachments = uploadFiles("rq1_attachment")
		req.Rq2Attachments = uploadFiles("rq2_attachment")
		req.Rq3Attachments = uploadFiles("rq3_attachment")
		req.Rq4Attachments = uploadFiles("rq4_attachment")
		req.Rq5Attachments = uploadFiles("rq5_attachment")
		req.Rq6Attachments = uploadFiles("rq6_attachment")
		req.Rq7Attachments = uploadFiles("rq7_attachment")

		req.Cq1Attachments = uploadFiles("cq1_attachment")
		req.Cq2Attachments = uploadFiles("cq2_attachment")
		req.Cq3Attachments = uploadFiles("cq3_attachment")
		req.Cq4Attachments = uploadFiles("cq4_attachment")
		req.Cq5Attachments = uploadFiles("cq5_attachment")
		req.Cq6Attachments = uploadFiles("cq6_attachment")
		req.Cq7Attachments = uploadFiles("cq7_attachment")
		req.Cq8Attachments = uploadFiles("cq8_attachment")
		req.Cq9Attachments = uploadFiles("cq9_attachment")

		if err := h.Repo.CreateAssessmentTrl(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
		return
	}

	// 2. Fallback to JSON
	var req models.AssessmentTrl
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAssessmentTrl(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /assessment/:id
func (h *AssessmentTrlHandler) UpdateAssessmentTrlByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateAssessmentTrlByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment TRL updated successfully"})
}
