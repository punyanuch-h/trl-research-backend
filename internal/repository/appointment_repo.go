package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type AppointmentRepo struct {
	DB *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) AppointmentRepository {
	return &AppointmentRepo{DB: db}
}

// 🟢 GetAppointmentAll
func (r *AppointmentRepo) GetAppointmentAll() ([]models.Appointments, error) {
	var appointments []models.Appointments
	err := r.DB.Find(&appointments).Error
	return appointments, err
}

// 🟢 GetAppointmentByID
func (r *AppointmentRepo) GetAppointmentByID(appointmentID string) (*models.Appointments, error) {
	var ap models.Appointments
	err := r.DB.Where("id = ?", appointmentID).First(&ap).Error
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

// 🟢 GetAppointmentByCaseID
func (r *AppointmentRepo) GetAppointmentByCaseID(caseID string) ([]models.Appointments, error) {
	var appointments []models.Appointments
	err := r.DB.Where("case_id = ?", caseID).Find(&appointments).Error
	return appointments, err
}

// 🟢 CreateAppointment - auto generate ID AP-00001
func (r *AppointmentRepo) CreateAppointment(ap *models.Appointments) error {
	id, err := utils.GenerateID(r.DB, "appointments", "AP")
	if err != nil {
		return err
	}
	ap.ID = id
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
	return r.DB.Model(&models.Appointments{}).Where("id = ?", appointmentID).Updates(data).Error
}
