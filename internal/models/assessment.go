package models

import (
	"time"

	"gorm.io/datatypes"
)

type Assessments struct {
	ID                    string         `gorm:"primaryKey;column:id" json:"id" form:"id"`
	CaseID                string         `gorm:"column:case_id;not null" json:"case_id" form:"case_id"`
	Case                  *Cases         `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case"`
	TrlEstimate           int16          `gorm:"column:trl_estimate;not null" json:"trl_estimate" form:"trl_estimate"`
	Rq1Answer             bool           `gorm:"column:rq1_answer;not null" json:"rq1_answer" form:"rq1_answer"`
	Rq1Attachments        datatypes.JSON `gorm:"column:rq1_attachments" json:"rq1_attachments" form:"-"`
	Rq2Answer             bool           `gorm:"column:rq2_answer;not null" json:"rq2_answer" form:"rq2_answer"`
	Rq2Attachments        datatypes.JSON `gorm:"column:rq2_attachments" json:"rq2_attachments" form:"-"`
	Rq3Answer             bool           `gorm:"column:rq3_answer;not null" json:"rq3_answer" form:"rq3_answer"`
	Rq3Attachments        datatypes.JSON `gorm:"column:rq3_attachments" json:"rq3_attachments" form:"-"`
	Rq4Answer             bool           `gorm:"column:rq4_answer;not null" json:"rq4_answer" form:"rq4_answer"`
	Rq4Attachments        datatypes.JSON `gorm:"column:rq4_attachments" json:"rq4_attachments" form:"-"`
	Rq5Answer             bool           `gorm:"column:rq5_answer;not null" json:"rq5_answer" form:"rq5_answer"`
	Rq5Attachments        datatypes.JSON `gorm:"column:rq5_attachments" json:"rq5_attachments" form:"-"`
	Rq6Answer             bool           `gorm:"column:rq6_answer;not null" json:"rq6_answer" form:"rq6_answer"`
	Rq6Attachments        datatypes.JSON `gorm:"column:rq6_attachments" json:"rq6_attachments" form:"-"`
	Rq7Answer             bool           `gorm:"column:rq7_answer;not null" json:"rq7_answer" form:"rq7_answer"`
	Rq7Attachments        datatypes.JSON `gorm:"column:rq7_attachments" json:"rq7_attachments" form:"-"`
	Cq1Answer             datatypes.JSON `gorm:"column:cq1_answer" json:"cq1_answer" form:"cq1_answer"`
	Cq1Attachments        datatypes.JSON `gorm:"column:cq1_attachments" json:"cq1_attachments" form:"-"`
	Cq2Answer             datatypes.JSON `gorm:"column:cq2_answer" json:"cq2_answer" form:"cq2_answer"`
	Cq2Attachments        datatypes.JSON `gorm:"column:cq2_attachments" json:"cq2_attachments" form:"-"`
	Cq3Answer             datatypes.JSON `gorm:"column:cq3_answer" json:"cq3_answer" form:"cq3_answer"`
	Cq3Attachments        datatypes.JSON `gorm:"column:cq3_attachments" json:"cq3_attachments" form:"-"`
	Cq4Answer             datatypes.JSON `gorm:"column:cq4_answer" json:"cq4_answer" form:"cq4_answer"`
	Cq4Attachments        datatypes.JSON `gorm:"column:cq4_attachments" json:"cq4_attachments" form:"-"`
	Cq5Answer             datatypes.JSON `gorm:"column:cq5_answer" json:"cq5_answer" form:"cq5_answer"`
	Cq5Attachments        datatypes.JSON `gorm:"column:cq5_attachments" json:"cq5_attachments" form:"-"`
	Cq6Answer             datatypes.JSON `gorm:"column:cq6_answer" json:"cq6_answer" form:"cq6_answer"`
	Cq6Attachments        datatypes.JSON `gorm:"column:cq6_attachments" json:"cq6_attachments" form:"-"`
	Cq7Answer             datatypes.JSON `gorm:"column:cq7_answer" json:"cq7_answer" form:"cq7_answer"`
	Cq7Attachments        datatypes.JSON `gorm:"column:cq7_attachments" json:"cq7_attachments" form:"-"`
	Cq8Answer             datatypes.JSON `gorm:"column:cq8_answer" json:"cq8_answer" form:"cq8_answer"`
	Cq8Attachments        datatypes.JSON `gorm:"column:cq8_attachments" json:"cq8_attachments" form:"-"`
	Cq9Answer             datatypes.JSON `gorm:"column:cq9_answer" json:"cq9_answer" form:"cq9_answer"`
	Cq9Attachments        datatypes.JSON `gorm:"column:cq9_attachments" json:"cq9_attachments" form:"-"`
	ImprovementSuggestion string         `gorm:"column:improvement_suggestion" json:"improvement_suggestion" form:"improvement_suggestion"`
	CreatedAt             time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Assessments) TableName() string {
	return "assessments"
}
