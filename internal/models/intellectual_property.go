package models

import (
	"time"

	"gorm.io/datatypes"
)

type IntellectualProperties struct {
	ID               string         `gorm:"primaryKey;column:id" json:"id" form:"id"`
	CaseID           string         `gorm:"column:case_id;not null" json:"case_id" form:"case_id"`
	Case 		     *Cases         `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case"`
	Types            string         `gorm:"column:types;not null" json:"types" form:"types"`
	ProtectionStatus string         `gorm:"column:protection_status;not null" json:"protection_status" form:"protection_status"`
	RequestNumber    string         `gorm:"column:request_number" json:"request_number" form:"request_number"`
	Attachments      datatypes.JSON `gorm:"column:attachments" json:"attachments" form:"-"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (IntellectualProperties) TableName() string {
	return "intellectual_properties"
}
