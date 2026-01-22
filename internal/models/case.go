package models

import (
	"time"

	"gorm.io/datatypes"
)

type CaseInfo struct {
	CaseID           string `gorm:"primaryKey;column:case_id" json:"case_id" form:"case_id"`
	CoordinatorEmail string `gorm:"column:coordinator_email" json:"coordinator_email" form:"coordinator_email"`
	TrlScore         string `gorm:"column:trl_score" json:"trl_score" form:"trl_score"`
	TrlSuggestion    string `gorm:"column:trl_suggestion" json:"trl_suggestion" form:"trl_suggestion"`
	Status           bool   `gorm:"column:status" json:"status" form:"status"`
	IsUrgent         bool   `gorm:"column:is_urgent" json:"is_urgent" form:"is_urgent"`
	UrgentReason     string `gorm:"column:urgent_reason" json:"urgent_reason" form:"urgent_reason"`
	UrgentFeedback   string `gorm:"column:urgent_feedback" json:"urgent_feedback" form:"urgent_feedback"`

	CaseTitle       string         `gorm:"column:case_title" json:"case_title" form:"case_title"`
	CaseType        string         `gorm:"column:case_type" json:"case_type" form:"case_type"`
	CaseDescription string         `gorm:"column:case_description" json:"case_description" form:"case_description"`
	CaseKeywords    string         `gorm:"column:case_keywords" json:"case_keywords" form:"case_keywords"`
	CaseAttachments datatypes.JSON `gorm:"column:case_attachments" json:"case_attachments" form:"-"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`

	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	ResearcherID string    `gorm:"column:researcher_id" json:"researcher_id" form:"researcher_id"`
}

func (CaseInfo) TableName() string {
	return "cases"
}
