package models

import (
	"time"

	"gorm.io/datatypes"
)

type Cases struct {
	ID             string         `gorm:"primaryKey;column:id" json:"id" form:"id"`
	ResearcherID   string         `gorm:"column:researcher_id;not null" json:"researcher_id" form:"researcher_id"`
	Researcher     *Researchers   `gorm:"foreignKey:ResearcherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"researcher"`
	CoordinatorID  string         `gorm:"column:coordinator_id;not null" json:"coordinator_id" form:"coordinator_id"`
	Coordinator    *Coordinators  `gorm:"foreignKey:CoordinatorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"coordinator"`
	Title          string         `gorm:"column:title;not null" json:"title" form:"title"`
	Type           string         `gorm:"column:type;not null" json:"type" form:"type"`
	Description    string         `gorm:"column:description;not null" json:"description" form:"description"`
	Keywords       string         `gorm:"column:keywords;not null" json:"keywords" form:"keywords"`
	Attachments    datatypes.JSON `gorm:"column:attachments" json:"attachments" form:"-"`
	TrlScore       *int16         `gorm:"column:trl_score" json:"trl_score" form:"trl_score"`
	TrlSuggestion  string         `gorm:"column:trl_suggestion" json:"trl_suggestion" form:"trl_suggestion"`
	Status         bool           `gorm:"column:status;not null" json:"status" form:"status"`
	IsUrgent       bool           `gorm:"column:is_urgent;not null" json:"is_urgent" form:"is_urgent"`
	UrgentReason   string         `gorm:"column:urgent_reason" json:"urgent_reason" form:"urgent_reason"`
	UrgentFeedback string         `gorm:"column:urgent_feedback" json:"urgent_feedback" form:"urgent_feedback"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Cases) TableName() string {
	return "cases"
}
