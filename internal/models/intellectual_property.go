package models

import (
	"time"
)

type IntellectualProperty struct {
	ID                 string    `json:"id" firestore:"id"`
	CaseID             string    `json:"case_id" firestore:"case_id"`
	IPTypes            string    `json:"ip_types" firestore:"ip_types"`
	IPProtectionStatus string    `json:"ip_protection_status" firestore:"ip_protection_status"`
	IPRequestNumber    string    `json:"ip_request_number" firestore:"ip_request_number"`
	CreatedAt          time.Time `json:"created_at" firestore:"created_at"`
}
