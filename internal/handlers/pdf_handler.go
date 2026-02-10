package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"trl-research-backend/internal/config"

	"github.com/gin-gonic/gin"
)

type PDFHandler struct {
	Cfg config.Config
}

type PDFRequest struct {
	HTML string `json:"html" binding:"required"`
}

func (h *PDFHandler) GeneratePDF(c *gin.Context) {
	var req PDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "HTML content is required"})
		return
	}

	// Forward request to Node.js PDF service
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	resp, err := http.Post(h.Cfg.PDFServiceURL+"/generate-pdf", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "PDF generation service is unavailable"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": "PDF service error", "details": string(body)})
		return
	}

	// Stream the PDF response back to the client
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.Error(err)
	}
}
