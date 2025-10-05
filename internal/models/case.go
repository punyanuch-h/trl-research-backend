package models

import (
	"time"
)

type CaseInfo struct {
	CaseID           string    `json:"case_id" firestore:"case_id"`
	ResearcherID     string    `json:"researcher_id" firestore:"researcher_id"`
	CoordinatorEmail string    `json:"coordinator_email" firestore:"coordinator_email"`
	TrlScore         string    `json:"trl_score" firestore:"trl_score"`
	Status           bool      `json:"status" firestore:"status"`
	IsUrgent         bool      `json:"is_urgent" firestore:"is_urgent"`
	UrgentReason     string    `json:"urgent_reason" firestore:"urgent_reason"`
	UrgentFeedback   string    `json:"urgent_feedback" firestore:"urgent_feedback"`
	CaseTitle        string    `json:"case_title" firestore:"case_title"`
	CaseType         string    `json:"case_type" firestore:"case_type"`
	CaseDescription  string    `json:"case_description" firestore:"case_description"`
	CaseKeywords     string    `json:"case_keywords" firestore:"case_keywords"`
	CreatedAt        time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" firestore:"updated_at"`
}
