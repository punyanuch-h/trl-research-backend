package models

import (
	"time"

	"gorm.io/datatypes"
)

type IntellectualProperty struct {
	ID                 string         `gorm:"primaryKey;column:id" json:"id" form:"id"`
	CaseID             string         `gorm:"column:case_id" json:"case_id" form:"case_id"`
	IPTypes            string         `gorm:"column:ip_types" json:"ip_types" form:"ip_types"`
	IPProtectionStatus string         `gorm:"column:ip_protection_status" json:"ip_protection_status" form:"ip_protection_status"`
	IPRequestNumber    string         `gorm:"column:ip_request_number" json:"ip_request_number" form:"ip_request_number"`
	IPAttachments      datatypes.JSON `gorm:"column:ip_attachment" json:"ip_attachment"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (IntellectualProperty) TableName() string {
	return "intellectual_properties"
}
