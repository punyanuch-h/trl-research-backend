package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils/send_email"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppointmentHandler struct {
	Repo repository.AppointmentRepository
	Cfg  config.Config
}

// 🟢 GET /appointments
func (h *AppointmentHandler) GetAppointmentAll(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	var appointments []models.Appointments
	var err error

	if role == "admin" {
		appointments, err = h.Repo.GetAppointmentAll()
	} else if role == "researcher" {
		uid, ok := userID.(string)
		if !ok || uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid user ID"})
			return
		}
		appointments, err = h.Repo.GetAppointmentByResearcherID(uid)
	} else {
		// Other roles see no appointments by default
		appointments = []models.Appointments{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// 🟢 GET /appointment/:id
func (h *AppointmentHandler) GetAppointmentByID(c *gin.Context) {
	id := c.Param("id")
	ap, err := h.Repo.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}
	c.JSON(http.StatusOK, ap)
}

// 🟢 GET /appointment/case/:id
func (h *AppointmentHandler) GetAppointmentByCaseID(c *gin.Context) {
	id := c.Param("id")
	appointments, err := h.Repo.GetAppointmentByCaseID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointments not found"})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

// 🟢 POST /appointment
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req models.Appointments
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAppointment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 📧 Send email notification
	go func(appointmentID string) {
		fullAp, err := h.Repo.GetAppointmentWithDetails(appointmentID)
		if err != nil {
			log.Printf("Failed to fetch appointment details for email: %v", err)
			return
		}

		details, _ := send_email.GatherAppointmentEmailData(fullAp)
		emailService := send_email.CreateSMTPEmailService(h.Cfg.GetSMTPConfig())
		results, err := send_email.SendEmail(details, emailService)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			log.Printf("Email send results: Sent=%d, Failed=%d, Skipped=%d", results.TotalSent, results.TotalFailed, results.TotalSkipped)
		}
	}(req.ID)

	c.JSON(http.StatusOK, req)
}

// 🟢 PATCH /appointment/:id
func (h *AppointmentHandler) UpdateAppointmentByID(c *gin.Context) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Fetch OLD appointment state for comparison
	oldAp, err := h.Repo.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	// 2. Logic for is_notify reset/set when date changes
	if newDateStr, ok := updateData["date"].(string); ok {
		// Attempt to parse date from common formats
		newDate, parseErr := time.Parse(time.RFC3339, newDateStr)
		if parseErr != nil {
			// Try RFC3339Nano as a secondary fallback
			newDate, parseErr = time.Parse(time.RFC3339Nano, newDateStr)
		}

		if parseErr == nil {
			now := time.Now()
			diff := newDate.Sub(now)

			if diff > time.Hour {
				// CASE A: Postponed to > 1 hour from now -> Reset is_notify to false
				updateData["is_notify"] = false
				log.Printf("[UpdateAppointment] Case A: Postponed > 1h. Resetting is_notify=false for %s", id)
			} else if diff >= 0 {
				// CASE B: Updated to start within 1 hour -> Set is_notify to true (update email acts as reminder)
				updateData["is_notify"] = true
				log.Printf("[UpdateAppointment] Case B: Starts within 1h. Setting is_notify=true for %s", id)
			}
		}
	} else if _, ok := updateData["is_notify"]; !ok {
		// If date didn't change, but other fields did and we are within 1h,
		// the update email will act as a reminder, so we can mark as notified.
		now := time.Now()
		if oldAp.Date.Sub(now) < time.Hour && oldAp.Date.After(now) {
			updateData["is_notify"] = true
			log.Printf("[UpdateAppointment] Meaningful update within 1h. Setting is_notify=true for %s", id)
		}
	}

	// 2.1 Always update 'updated_at' when an admin explicitly updates the appointment
	updateData["updated_at"] = time.Now()

	if err := h.Repo.UpdateAppointmentByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Fetch NEW appointment state
	newAp, err := h.Repo.GetAppointmentWithDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated appointment"})
		return
	}

	// 4. Trigger Update Notification logic (background)
	go h.ProcessUpdateNotification(oldAp, newAp)

	c.JSON(http.StatusOK, newAp)
}

// ProcessUpdateNotification handles meaningful change detection, anti-spam, and email sending
func (h *AppointmentHandler) ProcessUpdateNotification(oldAp *models.Appointments, newAp *models.Appointments) {
	// A. Detect meaningful changes
	changedFields := make(map[string]bool)
	meaningful := false

	if !newAp.Date.Equal(oldAp.Date) {
		changedFields["date"] = true
		meaningful = true
	}
	if newAp.Location != oldAp.Location {
		changedFields["location"] = true
		meaningful = true
	}
	if newAp.Detail != oldAp.Detail {
		changedFields["detail"] = true
		meaningful = true
	}
	if newAp.Summary != oldAp.Summary {
		changedFields["summary"] = true
		meaningful = true
	}

	if !meaningful {
		log.Printf("[Notification] No meaningful changes detected for appointment %s. Skipping email. Differences: %+v", newAp.ID, changedFields)
		return
	}

	log.Printf("[Notification] Meaningful changes detected for %s: %+v", newAp.ID, changedFields)

	// B. Anti-spam Protection
	now := time.Now()
	diff := now.Sub(oldAp.UpdatedAt)
	cooldown := h.Cfg.GetAntiSpamCooldown()
	if diff < cooldown {
		log.Printf("[Notification] Anti-spam: Skipping email for %s. Last meaningful update was only %v ago (limit: %v)", newAp.ID, diff.Round(time.Second), cooldown)
		return
	}
	log.Printf("[Notification] Anti-spam cleared for %s. Last update: %v (diff: %v, limit: %v)", newAp.ID, oldAp.UpdatedAt.Format("15:04:05"), diff.Round(time.Minute), cooldown)

	// C. Prepare Email Data
	details, _ := send_email.GatherAppointmentEmailData(newAp)
	emailService := send_email.CreateSMTPEmailService(h.Cfg.GetSMTPConfig())

	// Unique Recipients: Researcher and Coordinator
	recipients := make(map[string]string)
	if newAp.Case != nil {
		if newAp.Case.Researcher != nil {
			r := newAp.Case.Researcher
			recipients[r.Email] = fmt.Sprintf("%s %s %s", r.Prefix, r.FirstName, r.LastName)
		}
		if newAp.Case.Coordinator != nil {
			co := newAp.Case.Coordinator
			recipients[co.Email] = fmt.Sprintf("%s %s %s", co.Prefix, co.FirstName, co.LastName)
		}
	}

	if len(recipients) == 0 {
		log.Printf("[Notification] No recipients found for appointment %s (Case is nil: %v)", newAp.ID, newAp.Case == nil)
		return
	}
	log.Printf("[Notification] Found %d recipients for %s: %v", len(recipients), newAp.ID, recipients)

	subject := "[URGENT] Appointment Details Updated"
	updatedAtStr := now.Format("02 Jan 2006 15:04:05")

	// D. Send Emails
	anySuccess := false
	for email, name := range recipients {
		recipient := send_email.Recipient{Name: name, Email: email}
		content := send_email.TemplateUpdate(recipient, details, newAp.ID, updatedAtStr, changedFields)

		err := emailService(email, subject, content)
		if err != nil {
			log.Printf("❌ [Notification] Failed to send update email to %s: %v", email, err)
		} else {
			log.Printf("📧 [Notification] Sent update email to %s", email)
			anySuccess = true
		}
	}

	// E. Log results
	if anySuccess {
		log.Printf("[Notification] Successfully sent update emails for %s", newAp.ID)
	}
}

// 🟢 GET /trl/notifications/appointments
func (h *AppointmentHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	notifications, err := h.Repo.GetNotificationsByRole(role, userID)
	if err != nil {
		log.Printf("❌ [GetNotifications] Error fetching notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	unreadCount, err := h.Repo.GetUnreadNotificationCountByRole(role, userID)
	if err != nil {
		log.Printf("❌ [GetNotifications] Error fetching notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unread_count": unreadCount,
		"data":         notifications,
	})
}

// 🟢 PATCH /trl/notifications/appointments/:id/read
func (h *AppointmentHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ap, err := h.Repo.MarkNotificationAsRead(id, role, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found or access denied"})
			return
		}
		log.Printf("❌ [MarkAsRead] Error fetching notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, ap)
}

// 🟢 PATCH /trl/notifications/appointments/read-all
func (h *AppointmentHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.Repo.MarkAllNotificationsAsRead(role, userID); err != nil {
		log.Printf("❌ [MarkAllAsRead] Error fetching notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}
