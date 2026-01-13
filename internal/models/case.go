package models

import (
	"time"
)

type CaseInfo struct {
	CaseID           string    `json:"case_id" firestore:"case_id" form:"case_id"`
	CoordinatorEmail string    `json:"coordinator_email" firestore:"coordinator_email" form:"coordinator_email"`
	TrlScore         string    `json:"trl_score" firestore:"tr_score" form:"trl_score"`
	TrlSuggestion    string    `json:"trl_suggestion" firestore:"trl_suggestion" form:"trl_suggestion"`
	Status           bool      `json:"status" firestore:"status" form:"status"`
	IsUrgent         bool      `json:"is_urgent" firestore:"is_urgent" form:"is_urgent"`
	UrgentReason     string    `json:"urgent_reason" firestore:"urgent_reason" form:"urgent_reason"`
	UrgentFeedback   string    `json:"urgent_feedback" firestore:"urgent_feedback" form:"urgent_feedback"`
	CaseTitle        string    `json:"case_title" firestore:"case_title" form:"case_title"`
	CaseType         string    `json:"case_type" firestore:"case_type" form:"case_type"`
	CaseDescription  string    `json:"case_description" firestore:"case_description" form:"case_description"`
	CaseKeywords     string    `json:"case_keywords" firestore:"case_keywords" form:"case_keywords"`
	CaseAttachments  []string  `json:"case_attachments" firestore:"case_attachments"`
	CreatedAt        time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" firestore:"updated_at"`

	ResearcherID string `json:"researcher_id" firestore:"researcher_id" form:"researcher_id"`
}
