package models

import "time"

type Appointments struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	CaseID    string    `gorm:"column:case_id;not null" json:"case_id"`
	Case      *Cases    `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case"`
	Date      time.Time `gorm:"column:date;not null" json:"date"`
	Status    string    `gorm:"column:status;not null" json:"status"`
	Location  string    `gorm:"column:location" json:"location"`
	Detail    string    `gorm:"column:detail" json:"detail"`
	Summary   string    `gorm:"column:summary" json:"summary"`
	IsNotify  bool      `gorm:"column:is_notify" json:"is_notify"`
	IsRead    bool      `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Appointments) TableName() string {
	return "appointments"
}
