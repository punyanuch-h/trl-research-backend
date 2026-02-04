package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type IntellectualPropertyHandler struct {
	Repo repository.IntellectualPropertyRepository
	GCS  *storage.GCSClient
}

// 🟢 GET /ips
func (h *IntellectualPropertyHandler) GetIPAll(c *gin.Context) {
	ips, err := h.Repo.GetIPAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ips)
}

// 🟢 GET /ip/:id
func (h *IntellectualPropertyHandler) GetIPByID(c *gin.Context) {
	id := c.Param("id")
	ip, err := h.Repo.GetIPByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Intellectual Property not found"})
		return
	}
	c.JSON(http.StatusOK, ip)
}

// 🟢 GET /ip/case/:id
func (h *IntellectualPropertyHandler) GetIPByCaseID(c *gin.Context) {
	id := c.Param("id")
	ip, err := h.Repo.GetIPByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Intellectual Property not found"})
		return
	}
	c.JSON(http.StatusOK, ip)
}

// 🟢 POST /ip
func (h *IntellectualPropertyHandler) CreateIP(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")

	var req models.IntellectualProperties
	var attachments []string

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (strings.Contains(contentType, "multipart/form-data")) {
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		files := form.File["ips_attachments"]

		today := time.Now().Format("2006-01-02")
		userID := c.GetString("userID")
		if userID == "" {
			userID = "unknown_user"
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				continue
			}
			defer file.Close()

			objectPath := fmt.Sprintf("attachments/ips/%s/%s/%s", userID, today, fileHeader.Filename)
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
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}
		if err := json.Unmarshal(bodyJSON, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		// Extract semantic attachments using utility
		attachmentsMap := utils.ExtractAttachments(body)
		if paths, ok := attachmentsMap["ips"]; ok {
			attachments = paths
		}
	}

	// Save attachments if any
	if len(attachments) > 0 {
		jsonData, err := json.Marshal(attachments)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse attachments"})
			return
		}
		req.Attachments = datatypes.JSON(jsonData)
	}

	if err := h.Repo.CreateIP(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /ip/:id
func (h *IntellectualPropertyHandler) UpdateIPByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract semantic attachments using utility
	attachmentsMap := utils.ExtractAttachments(updateData)
	_, ipKeyPresent := updateData["ips_attachments"]

	if paths, ok := attachmentsMap["ips"]; ok {
		jsonData, _ := json.Marshal(paths)
		updateData["attachments"] = string(jsonData)
	} else if ipKeyPresent {
		updateData["attachments"] = "[]"
	}

	// Always remove semantic keys to avoid passing them to repository/DB
	for key := range utils.AttachmentKeys {
		delete(updateData, key)
	}

	if err := h.Repo.UpdateIPByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Intellectual Property updated successfully"})
}
