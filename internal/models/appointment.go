package models

import "time"

type Appointment struct {
	AppointmentID string    `gorm:"primaryKey;column:appointment_id" json:"appointment_id"`
	CaseID        string    `gorm:"column:case_id" json:"case_id"`
	Date          time.Time `gorm:"column:date" json:"date"`
	Status        string    `gorm:"column:status" json:"status"`
	Location      string    `gorm:"column:location" json:"location"`
	Note          string    `gorm:"column:note" json:"note"`
	Summary       string    `gorm:"column:summary" json:"summary"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Appointment) TableName() string {
	return "appointments"
}
