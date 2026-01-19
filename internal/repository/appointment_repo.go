package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type AppointmentRepo struct {
	DB *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) AppointmentRepository {
	return &AppointmentRepo{DB: db}
}

// 🟢 GetAppointmentAll
func (r *AppointmentRepo) GetAppointmentAll() ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Find(&appointments).Error
	return appointments, err
}

// 🟢 GetAppointmentByID
func (r *AppointmentRepo) GetAppointmentByID(appointmentID string) (*models.Appointment, error) {
	var ap models.Appointment
	err := r.DB.Where("appointment_id = ?", appointmentID).First(&ap).Error
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

// 🟢 GetAppointmentByCaseID
func (r *AppointmentRepo) GetAppointmentByCaseID(caseID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Where("case_id = ?", caseID).Find(&appointments).Error
	return appointments, err
}

// 🟢 CreateAppointment - auto generate ID AP-00001
func (r *AppointmentRepo) CreateAppointment(ap *models.Appointment) error {
	var lastAp models.Appointment
	nextID := "AP-00001"
	err := r.DB.Order("appointment_id desc").First(&lastAp).Error
	if err == nil {
		lastID := lastAp.AppointmentID
		numStr := strings.TrimPrefix(lastID, "AP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("AP-%05d", n+1)
		}
	}

	ap.AppointmentID = nextID
	now := time.Now()
	ap.CreatedAt = now
	ap.UpdatedAt = now

	return r.DB.Create(ap).Error
}

// 🟢 UpdateAppointmentByID
func (r *AppointmentRepo) UpdateAppointmentByID(appointmentID string, data map[string]interface{}) error {
	// Convert date string to time.Time if present
	if dateStr, ok := data["date"].(string); ok {
		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			// Try alternative format if RFC3339 fails
			parsedDate, err = time.Parse("2006-01-02T15:04:05Z07:00", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format: %v", err)
			}
		}
		data["date"] = parsedDate
	}

	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Appointment{}).Where("appointment_id = ?", appointmentID).Updates(data).Error
}
