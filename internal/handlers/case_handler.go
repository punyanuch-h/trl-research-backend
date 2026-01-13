package handlers

import (
	"fmt"
	"net/http"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type CaseHandler struct {
	Repo     *repository.CaseRepo
	FileRepo *repository.FileRepo
	GCS      *storage.GCSClient
}

// 🟢 GET /cases
func (h *CaseHandler) GetCaseAll(c *gin.Context) {
	fmt.Println("GetCaseAll from handler")
	fmt.Println("h", h)
	cases, err := h.Repo.GetCaseAll()
	fmt.Println("cases", cases)
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cases)
}

// 🟢 GET /case/researcher/:id - Get all cases for a researcher
func (h *CaseHandler) GetCaseAllByResearcher_id(c *gin.Context) {
	id := c.Param("id")
	cases, err := h.Repo.GetCaseAllByResearcher_id(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cases)
}

// 🟢 GET /case/:id
func (h *CaseHandler) GetCaseByID(c *gin.Context) {
	id := c.Param("id")
	cs, err := h.Repo.GetCaseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}
	c.JSON(http.StatusOK, cs)
}

// 🟢 POST /case
func (h *CaseHandler) CreateCase(c *gin.Context) {
	// Check Content-Type to determine how to parse
	contentType := c.GetHeader("Content-Type")

	// 1. Handle Multipart/Form-Data (File Upload)
	if contentType != "" && (contentType == "multipart/form-data" || len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		var req models.CaseInfo
		// Bind form fields to struct
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		// Handle Multiple Files Upload (key: "case_attachments")
		form, _ := c.MultipartForm()
		files := form.File["case_attachments"]

		var uploadedPaths []string

		userID := req.ResearcherID
		if userID == "" {
			userID = c.GetString("userID")
		}
		if userID == "" {
			userID = "unknown_user"
		}
		today := time.Now().Format("2006-01-02")

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file: " + err.Error()})
				return
			}
			defer file.Close()

			objectPath := fmt.Sprintf("case_attachments/%s/%s/%s", today, userID, fileHeader.Filename)

			// Upload to GCS
			if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file: " + err.Error()})
				return
			}
			uploadedPaths = append(uploadedPaths, objectPath)
		}

		// Add paths to request model
		req.CaseAttachments = uploadedPaths

		// Save Case
		if err := h.Repo.CreateCase(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create case: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
		return
	}

	// 2. Fallback to JSON (Existing Logic)
	var req models.CaseInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateCase(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /case/:id
func (h *CaseHandler) UpdateCaseByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.UpdateCaseByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Case updated successfully"})
}

// 🟢 PATCH /case/update-status/:id
func (h *CaseHandler) UpdateCaseStatusByID(c *gin.Context) {
	id := c.Param("id")
	status := c.Query("status")

	if err := h.Repo.UpdateCaseStatusByID(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Case status updated successfully"})
}
