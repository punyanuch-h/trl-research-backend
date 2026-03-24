package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"
	"trl-research-backend/internal/utils"
	"trl-research-backend/internal/utils/send_email"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type CaseHandler struct {
	Repo repository.CaseRepository
	GCS  *storage.GCSClient
	Cfg  config.Config
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
	log.Println("[CreateCase] Incoming request")
	contentType := c.GetHeader("Content-Type")

	var req models.Cases
	var attachments []string

	// 1. Handle Multipart/Form-Data
	if contentType != "" && (strings.Contains(contentType, "multipart/form-data")) {
		log.Println("[CreateCase] Handling multipart/form-data upload")
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("[CreateCase] Form bind error: %v\n", err)
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
				log.Printf("[CreateCase] Error opening file %s: %v\n", fileHeader.Filename, err)
				continue
			}
			defer file.Close()

			// Build path: attachments/cases/{userID}/{date}/{filename}
			objectPath := fmt.Sprintf("attachments/cases/%s/%s/%s", userID, today, fileHeader.Filename)
			log.Printf("[CreateCase] Uploading file to GCS: %s\n", objectPath)

			if err := h.GCS.UploadFile(objectPath, fileHeader.Header.Get("Content-Type"), file); err == nil {
				log.Printf("[CreateCase] Successfully uploaded %s\n", objectPath)
				attachments = append(attachments, objectPath)
			} else {
				log.Printf("[CreateCase] Failed to upload %s: %v\n", objectPath, err)
			}
		}
	} else {
		// 2. Handle JSON
		log.Println("[CreateCase] Handling JSON request (likely using pre-signed URL paths)")
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			log.Printf("[CreateCase] JSON bind error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Use dynamic mapping to populate the struct
		bodyJSON, _ := json.Marshal(body)
		if err := json.Unmarshal(bodyJSON, &req); err != nil {
			log.Printf("[CreateCase] JSON unmarshal error: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request body"})
			return
		}

		// Extract semantic attachments using utility
		attachmentsMap := utils.ExtractAttachments(body)
		if paths, ok := attachmentsMap["cases"]; ok {
			log.Printf("[CreateCase] Found %d attachments in JSON request\n", len(paths))
			attachments = paths
		}
	}

	// Save attachments if any
	if len(attachments) > 0 {
		log.Printf("[CreateCase] Saving %d attachments to DB record\n", len(attachments))
		jsonData, _ := json.Marshal(attachments)
		req.Attachments = datatypes.JSON(jsonData)
	}

	if err := h.Repo.CreateCase(&req); err != nil {
		log.Printf("[CreateCase] Repo creation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[CreateCase] Successfully created case ID: %s\n", req.ID)
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

	// 1. Fetch OLD case state for comparison
	oldCase, err := h.Repo.GetCaseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
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

	// 2. Fetch NEW case state
	newCase, err := h.Repo.GetCaseByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated case"})
		return
	}

	// 3. Trigger Notification logic (background)
	go h.ProcessCaseStatusNotification(oldCase, newCase)

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

	// 1. Fetch OLD case state for comparison
	oldCase, err := h.Repo.GetCaseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	if err := h.Repo.UpdateCaseStatusByID(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Fetch NEW case state
	newCase, err := h.Repo.GetCaseByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated case"})
		return
	}

	// 3. Trigger Notification logic (background)
	go h.ProcessCaseStatusNotification(oldCase, newCase)

	c.JSON(http.StatusOK, gin.H{"message": "Case status updated successfully"})
}

// ProcessCaseStatusNotification handles status change detection, anti-spam, and email sending
func (h *CaseHandler) ProcessCaseStatusNotification(oldCase *models.Cases, newCase *models.Cases) {
	// A. Detect meaningful change: Send email ONLY when status changes from FALSE to TRUE (Approved)
	statusChangedToApproved := !oldCase.Status && newCase.Status

	if !statusChangedToApproved {
		log.Printf("[CaseNotification] Skipping email for case %s. Status did not change to 'Approved'. (Old: %v, New: %v)", newCase.ID, oldCase.Status, newCase.Status)
		return
	}

	// If we reach here, it means the case was just approved.
	log.Printf("[CaseNotification] Status approved for %s: %v -> %v. Preparing to send email.", newCase.ID, oldCase.Status, newCase.Status)

	// B. Anti-spam Protection
	now := time.Now()

	// C. Prepare Email Data
	emailService := send_email.CreateSMTPEmailService(h.Cfg.GetSMTPConfig())

	// Unique Recipients: Researcher and Coordinator
	recipients := make(map[string]string)
	if newCase.Researcher != nil {
		r := newCase.Researcher
		recipients[r.Email] = fmt.Sprintf("%s %s %s", r.Prefix, r.FirstName, r.LastName)
	}
	if newCase.Coordinator != nil {
		co := newCase.Coordinator
		recipients[co.Email] = fmt.Sprintf("%s %s %s", co.Prefix, co.FirstName, co.LastName)
	}

	if len(recipients) == 0 {
		log.Printf("⚠️ [CaseNotification] No recipients found for Case %s. (Did you Preload Researcher/Coordinator in GetCaseByID?)", newCase.ID)
		return
	}
	log.Printf("[CaseNotification] Found %d recipients for %s: %v", len(recipients), newCase.ID, recipients)

	subject := "[Update] Research Project Status Changed"
	updatedAtStr := now.Format("02 Jan 2006 15:04:05")

	// Determine Project Title (fallback to ID if Title field is empty or doesn't exist)
	title := newCase.Title
	if title == "" {
		title = newCase.ID
	}

	// D. Send Emails
	anySuccess := false
	for email, name := range recipients {
		recipient := send_email.Recipient{Name: name, Email: email}
		content := send_email.TemplateCaseStatusChange(
			recipient, title, updatedAtStr,
		)

		err := emailService(email, subject, content)
		if err != nil {
			log.Printf("❌ [CaseNotification] Failed to send update email to %s: %v", email, err)
		} else {
			log.Printf("📧 [CaseNotification] Sent update email to %s", email)
			anySuccess = true
		}
	}

	// E. Log results
	if anySuccess {
		log.Printf("[CaseNotification] Successfully sent status update emails for %s", newCase.ID)
	}
}
