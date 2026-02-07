package handlers

import (
	"errors"
	"log"
	"net/http"

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

	if err := h.Repo.UpdateAppointmentByID(id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedAppointment, err := h.Repo.GetAppointmentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated appointment"})
		return
	}

	c.JSON(http.StatusOK, updatedAppointment)
}

// 🟢 GET /trl/notifications/appointments
func (h *AppointmentHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")
	if userID == "" || role == "" {
		log.Printf("❌ [GetNotifications] Error fetching notifications: %v", err)
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
