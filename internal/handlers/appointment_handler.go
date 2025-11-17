package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
)

type AppointmentHandler struct {
	Repo *repository.AppointmentRepo
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
	var req models.Appointment
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.CreateAppointment(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully"})
}
