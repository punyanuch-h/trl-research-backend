package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CaseHandler struct {
	Repo     repository.CaseRepository
	FileRepo repository.FileRepository
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
func (h *CaseHandler) GetCaseAllByResearcherID(c *gin.Context) {
	id := c.Param("id")
	cases, err := h.Repo.GetCaseAllByResearcherID(id)
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
	contentType := c.GetHeader("Content-Type")

	var req models.Cases
	var attachments []string

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (strings.Contains(contentType, "multipart/form-data")) {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		form, _ := c.MultipartForm()
		files := form.File["cases_attachments"]

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
				continue
			}
			defer file.Close()

			// Build path: attachments/cases/{userID}/{date}/{filename}
			objectPath := fmt.Sprintf("attachments/cases/%s/%s/%s", userID, today, fileHeader.Filename)

			if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err == nil {
				attachments = append(attachments, objectPath)
			}
		}
	} else {
		// 2. Handle JSON
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use dynamic mapping to populate the struct
		bodyJSON, _ := json.Marshal(body)
		if err := json.Unmarshal(bodyJSON, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		// Extract semantic attachments using utility
		attachmentsMap := utils.ExtractAttachments(body)
		if paths, ok := attachmentsMap["cases"]; ok {
			attachments = paths
		}
	}

	// Save attachments if any
	if len(attachments) > 0 {
		jsonData, _ := json.Marshal(attachments)
		req.Attachments = datatypes.JSON(jsonData)
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

	// Extract semantic attachments using utility
	attachmentsMap := utils.ExtractAttachments(updateData)
	_, caseKeyPresent := updateData["cases_attachments"]

	if paths, ok := attachmentsMap["cases"]; ok {
		jsonData, _ := json.Marshal(paths)
		updateData["attachments"] = string(jsonData)
	} else if caseKeyPresent {
		updateData["attachments"] = "[]"
	}

	// Always remove semantic keys to avoid passing them to repository/DB
	for key := range utils.AttachmentKeys {
		delete(updateData, key)
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
	statusStr := c.Query("status")

	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value, expected boolean"})
		return
	}

	if err := h.Repo.UpdateCaseStatusByID(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Case status updated successfully"})
}
