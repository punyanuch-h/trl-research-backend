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
	err := r.DB.Preload("Case").Find(&appointments).Error
	return appointments, err
}

// 🟢 GetAppointmentByID
func (r *AppointmentRepo) GetAppointmentByID(appointmentID string) (*models.Appointments, error) {
	var ap models.Appointments
	err := r.DB.Preload("Case").Where("id = ?", appointmentID).First(&ap).Error
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

// 🟢 GetAppointmentWithDetails
func (r *AppointmentRepo) GetAppointmentWithDetails(appointmentID string) (*models.Appointments, error) {
	var ap models.Appointments
	err := r.DB.Preload("Case").Preload("Case.Researcher").Preload("Case.Coordinator").Where("id = ?", appointmentID).First(&ap).Error
	if err != nil {
		return nil, err
	}
	return &ap, nil
}

// 🟢 GetAppointmentByCaseID
func (r *AppointmentRepo) GetAppointmentByCaseID(caseID string) ([]models.Appointments, error) {
	var appointments []models.Appointments
	err := r.DB.Preload("Case").Where("case_id = ?", caseID).Find(&appointments).Error
	return appointments, err
}

// 🟢 GetAppointmentByResearcherID
func (r *AppointmentRepo) GetAppointmentByResearcherID(researcherID string) ([]models.Appointments, error) {
	var appointments []models.Appointments
	err := r.DB.Joins("JOIN cases ON cases.id = appointments.case_id").
		Select("appointments.*").
		Preload("Case").
		Where("cases.researcher_id = ?", researcherID).
		Find(&appointments).Error
	return appointments, err
}

// 🟢 CreateAppointment - auto generate ID AP-00001
func (r *AppointmentRepo) CreateAppointment(ap *models.Appointments) error {
	id, err := utils.GenerateID(r.DB, "appointments", "AP")
	if err != nil {
		return err
	}
	ap.ID = id
	ap.IsNotify = false
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

// 🟢 GetUpcomingAppointments
func (r *AppointmentRepo) GetUpcomingAppointments(start, end time.Time) ([]models.Appointments, error) {
	var appointments []models.Appointments
	err := r.DB.Preload("Case").
		Preload("Case.Researcher").
		Preload("Case.Coordinator").
		Where("date BETWEEN ? AND ?", start, end).
		Where("is_notify = ?", false).
		Find(&appointments).Error
	return appointments, err
}

// 🟢 UpdateNotifyStatus
func (r *AppointmentRepo) UpdateNotifyStatus(appointmentID string, isNotify bool) error {
	return r.DB.Model(&models.Appointments{}).
		Where("id = ?", appointmentID).
		Update("is_notify", isNotify).Error
}

// 🟢 GetNotificationsByRole
func (r *AppointmentRepo) GetNotificationsByRole(role string, userID string) ([]models.Appointments, error) {
	var appointments []models.Appointments
	query := r.DB.Preload("Case").Where("is_notify = ?", true)

	if role == "admin" {
		// Admin sees all
	} else if role == "researcher" {
		query = query.Joins("JOIN cases ON cases.id = appointments.case_id").
			Where("cases.researcher_id = ?", userID)
	} else {
		// Other roles see no notifications by default
		return []models.Appointments{}, nil
	}

	err := query.Select("appointments.*").
		Order("is_read ASC").    // unread first
		Order("created_at DESC"). // newest first
		Find(&appointments).Error
	return appointments, err
}

// 🟢 GetUnreadNotificationCountByRole
func (r *AppointmentRepo) GetUnreadNotificationCountByRole(role string, userID string) (int64, error) {
	var count int64
	query := r.DB.Model(&models.Appointments{}).
		Where("is_notify = ? AND is_read = ?", true, false)

	if role == "admin" {
		// Admin sees all
	} else if role == "researcher" {
		query = query.Joins("JOIN cases ON cases.id = appointments.case_id").
			Where("cases.researcher_id = ?", userID)
	} else {
		return 0, nil
	}

	err := query.Count(&count).Error
	return count, err
}

// 🟢 MarkNotificationAsRead
func (r *AppointmentRepo) MarkNotificationAsRead(id string, role string, userID string) (*models.Appointments, error) {
	var ap models.Appointments
	query := r.DB.Where("appointments.id = ? AND is_notify = ?", id, true)
	
	if role == "admin" {
		// Admin can mark any as read
	} else if role == "researcher" {
		query = query.Joins("JOIN cases ON cases.id = appointments.case_id").
			Where("cases.researcher_id = ?", userID)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	
	err := query.Select("appointments.*").First(&ap).Error
	if err != nil {
		return nil, err
	}
	
	err = r.DB.Model(&ap).Update("is_read", true).Error
	if err != nil {
		return nil, err
	}
	
	ap.IsRead = true
	return &ap, nil
}

// 🟢 MarkAllNotificationsAsRead
func (r *AppointmentRepo) MarkAllNotificationsAsRead(role string, userID string) error {
	if role == "admin" {
		return r.DB.Model(&models.Appointments{}).
			Where("is_notify = ?", true).
			Update("is_read", true).Error
	} else if role == "researcher" {
		// Mark all notifications as read for appointments owned by this researcher
		subQuery := r.DB.Table("appointments").
			Select("appointments.id").
			Joins("JOIN cases ON cases.id = appointments.case_id").
			Where("cases.researcher_id = ?", userID).
			Where("appointments.is_notify = ?", true)
		
		return r.DB.Model(&models.Appointments{}).
			Where("id IN (?)", subQuery).
			Update("is_read", true).Error
	}

	return nil
}
