package handlers

import (
	"log"
	"net/http"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils/send_email"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	Repo repository.AppointmentRepository
	Cfg  config.Config
}

// 🟢 GET /appointments
func (h *AppointmentHandler) GetAppointmentAll(c *gin.Context) {
	appointments, err := h.Repo.GetAppointmentAll()
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
		smtpConfig := send_email.SMTPConfig{
			Host:     h.Cfg.EmailHost,
			Port:     h.Cfg.EmailPort,
			Username: h.Cfg.EmailSender,
			Password: h.Cfg.EmailPassword,
			From:     h.Cfg.EmailSender,
		}

		emailService := send_email.CreateSMTPEmailService(smtpConfig)
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
