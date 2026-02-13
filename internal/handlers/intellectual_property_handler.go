package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("[CreateIP] Incoming request")
	contentType := c.GetHeader("Content-Type")

	var req models.IntellectualProperties
	var attachments []string

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (strings.Contains(contentType, "multipart/form-data")) {
		log.Println("[CreateIP] Handling multipart/form-data upload")
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("[CreateIP] Form bind error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			log.Printf("[CreateIP] Multipart form error: %v\n", err)
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
				log.Printf("[CreateIP] Error opening file %s: %v\n", fileHeader.Filename, err)
				continue
			}
			defer file.Close()

			objectPath := fmt.Sprintf("attachments/ips/%s/%s/%s", userID, today, fileHeader.Filename)
			log.Printf("[CreateIP] Uploading file to GCS: %s\n", objectPath)
			if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err == nil {
				log.Printf("[CreateIP] Successfully uploaded %s\n", objectPath)
				attachments = append(attachments, objectPath)
			} else {
				log.Printf("[CreateIP] Failed to upload %s: %v\n", objectPath, err)
			}
		}
	} else {
		// 2. Handle JSON
		log.Println("[CreateIP] Handling JSON request (likely using pre-signed URL paths)")
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			log.Printf("[CreateIP] JSON bind error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use dynamic mapping to populate the struct
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			log.Printf("[CreateIP] JSON marshal error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}
		if err := json.Unmarshal(bodyJSON, &req); err != nil {
			log.Printf("[CreateIP] JSON unmarshal error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		// Extract semantic attachments using utility
		attachmentsMap := utils.ExtractAttachments(body)
		if paths, ok := attachmentsMap["ips"]; ok {
			log.Printf("[CreateIP] Found %d attachments in JSON request\n", len(paths))
			attachments = paths
		}
	}

	// Save attachments if any
	if len(attachments) > 0 {
		log.Printf("[CreateIP] Saving %d attachments to DB record\n", len(attachments))
		jsonData, err := json.Marshal(attachments)
		if err != nil {
			log.Printf("[CreateIP] JSON marshal error for attachments: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse attachments"})
			return
		}
		req.Attachments = datatypes.JSON(jsonData)
	}

	if err := h.Repo.CreateIP(&req); err != nil {
		log.Printf("[CreateIP] Repo creation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[CreateIP] Successfully created IP record with ID: %s\n", req.ID)
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
