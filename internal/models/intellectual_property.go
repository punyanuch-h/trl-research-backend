package models

import (
	"time"
)

type IntellectualProperty struct {
	ID                 string   `json:"id" firestore:"id" form:"id"`
	CaseID             string   `json:"case_id" firestore:"case_id" form:"case_id"`
	IPTypes            string   `json:"ip_types" firestore:"ip_types" form:"ip_types"`
	IPProtectionStatus string   `json:"ip_protection_status" firestore:"ip_protection_status" form:"ip_protection_status"`
	IPRequestNumber    string   `json:"ip_request_number" firestore:"ip_request_number" form:"ip_request_number"`
	IPAttachments      []string `json:"ip_attachment" firestore:"ip_attachment"`

	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}
