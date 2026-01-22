package models

import (
	"time"
)

type CoordinatorInfo struct {
	CoordinatorID    string    `gorm:"primaryKey;column:coordinator_id" json:"coordinator_id"`
	CoordinatorEmail string    `gorm:"column:coordinator_email;unique" json:"coordinator_email"`
	CoordinatorName  string    `gorm:"column:coordinator_name" json:"coordinator_name"`
	CoordinatorPhone string    `gorm:"column:coordinator_phone" json:"coordinator_phone"`
	Department       string    `gorm:"column:department" json:"department"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	CaseID           string    `gorm:"column:case_id" json:"case_id"`
}

func (CoordinatorInfo) TableName() string {
	return "coordinators"
}
