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

type IntellectualPropertyHandler struct {
	Repo *repository.IntellectualPropertyRepo
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

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (len(contentType) > 19 && contentType[:19] == "multipart/form-data") {
		var req models.IntellectualProperty
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
			return
		}

		form, _ := c.MultipartForm()
		files := form.File["ip_attachment"]
		var uploadedPaths []string

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

			objectPath := fmt.Sprintf("ip_attachments/%s/%s/%s", today, userID, fileHeader.Filename)
			if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err == nil {
				uploadedPaths = append(uploadedPaths, objectPath)
			}
		}

		req.IPAttachments = uploadedPaths

		if err := h.Repo.CreateIP(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, req)
		return
	}

	// 2. Fallback to JSON
	var req models.IntellectualProperty
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	if err := h.Repo.UpdateIPByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Intellectual Property updated successfully"})
}
