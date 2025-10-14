package models

import (
	"time"
)

type CoordinatorInfo struct {
	CoordinatorID    string    `json:"coordinator_id" firestore:"coordinator_id"`
	CoordinatorEmail string    `json:"coordinator_email" firestore:"coordinator_email"`
	CoordinatorName  string    `json:"coordinator_name" firestore:"coordinator_name"`
	CoordinatorPhone string    `json:"coordinator_phone" firestore:"coordinator_phone"`
	Department       string    `json:"department" firestore:"department"`
	CreatedAt        time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" firestore:"updated_at"`

	CaseID           string    `json:"case_id" firestore:"case_id"`
}
