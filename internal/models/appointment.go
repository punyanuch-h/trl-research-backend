package models

import "time"

type Appointment struct {
	ID       string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CaseID   string    `gorm:"type:varchar(10);not null"`
	Date     time.Time `gorm:"not null"`
	Status   string    `gorm:"type:text;not null"`
	Location string    `gorm:"type:text;not null"`
	Summary  string    `gorm:"type:text"`
	Note     string    `gorm:"type:text"`
}
