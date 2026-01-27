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

type AssessmentHandler struct {
	Repo repository.AssessmentRepository
	GCS  *storage.GCSClient
}

// 🟢 GET /assessments
func (h *AssessmentHandler) GetAssessmentAll(c *gin.Context) {
	assessments, err := h.Repo.GetAssessmentAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assessments)
}

// 🟢 GET /assessment/:id
func (h *AssessmentHandler) GetAssessmentByID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.Repo.GetAssessmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// 🟢 GET /assessment/case/:id
func (h *AssessmentHandler) GetAssessmentByCaseID(c *gin.Context) {
	id := c.Param("id")
	a, err := h.Repo.GetAssessmentByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// 🟢 POST /assessment
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		var req models.Assessments
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

		// Helper to bind JSON fields from form (for answers)
		bindJSON := func(key string) datatypes.JSON {
			val := c.PostForm(key)
			if val == "" {
				return datatypes.JSON("[]")
			}

			// Check if it's a valid JSON array
			if !json.Valid([]byte(val)) {
				return datatypes.JSON("[]")
			}

			return datatypes.JSON(val)
		}

		// Process all attachments
		req.Rq1Attachments = uploadFiles("rq1_attachments")
		req.Rq2Attachments = uploadFiles("rq2_attachments")
		req.Rq3Attachments = uploadFiles("rq3_attachments")
		req.Rq4Attachments = uploadFiles("rq4_attachments")
		req.Rq5Attachments = uploadFiles("rq5_attachments")
		req.Rq6Attachments = uploadFiles("rq6_attachments")
		req.Rq7Attachments = uploadFiles("rq7_attachments")

		// Process all answers
		req.Cq1Answer = bindJSON("cq1_answer")
		req.Cq2Answer = bindJSON("cq2_answer")
		req.Cq3Answer = bindJSON("cq3_answer")
		req.Cq4Answer = bindJSON("cq4_answer")
		req.Cq5Answer = bindJSON("cq5_answer")
		req.Cq6Answer = bindJSON("cq6_answer")
		req.Cq7Answer = bindJSON("cq7_answer")
		req.Cq8Answer = bindJSON("cq8_answer")
		req.Cq9Answer = bindJSON("cq9_answer")

		req.Cq1Attachments = uploadFiles("cq1_attachments")
		req.Cq2Attachments = uploadFiles("cq2_attachments")
		req.Cq3Attachments = uploadFiles("cq3_attachments")
		req.Cq4Attachments = uploadFiles("cq4_attachments")
		req.Cq5Attachments = uploadFiles("cq5_attachments")
		req.Cq6Attachments = uploadFiles("cq6_attachments")
		req.Cq7Attachments = uploadFiles("cq7_attachments")
		req.Cq8Attachments = uploadFiles("cq8_attachments")
		req.Cq9Attachments = uploadFiles("cq9_attachments")

		if err := h.Repo.CreateAssessment(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
		return
	}

	// 2. Fallback to JSON
	var req models.Assessments
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAssessment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /assessment/:id
func (h *AssessmentHandler) UpdateAssessmentByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateAssessmentByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment updated successfully"})
}
