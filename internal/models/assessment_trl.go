package models

import (
	"time"

	"gorm.io/datatypes"
)

type AssessmentTrl struct {
	ID             string         `gorm:"primaryKey;column:id" json:"id" form:"id"`
	CaseID         string         `gorm:"column:case_id" json:"case_id" form:"case_id"`
	TrlLevelResult int            `gorm:"column:trl_level_result" json:"trl_level_result" form:"trl_level_result"`
	Rq1Answer      bool           `gorm:"column:rq1_answer" json:"rq1_answer" form:"rq1_answer"`
	Rq1Attachments datatypes.JSON `gorm:"column:rq1_attachment" json:"rq1_attachment"`
	Rq2Answer      bool           `gorm:"column:rq2_answer" json:"rq2_answer" form:"rq2_answer"`
	Rq2Attachments datatypes.JSON `gorm:"column:rq2_attachment" json:"rq2_attachment"`
	Rq3Answer      bool           `gorm:"column:rq3_answer" json:"rq3_answer" form:"rq3_answer"`
	Rq3Attachments datatypes.JSON `gorm:"column:rq3_attachment" json:"rq3_attachment"`
	Rq4Answer      bool           `gorm:"column:rq4_answer" json:"rq4_answer" form:"rq4_answer"`
	Rq4Attachments datatypes.JSON `gorm:"column:rq4_attachment" json:"rq4_attachment"`
	Rq5Answer      bool           `gorm:"column:rq5_answer" json:"rq5_answer" form:"rq5_answer"`
	Rq5Attachments datatypes.JSON `gorm:"column:rq5_attachment" json:"rq5_attachment"`
	Rq6Answer      bool           `gorm:"column:rq6_answer" json:"rq6_answer" form:"rq6_answer"`
	Rq6Attachments datatypes.JSON `gorm:"column:rq6_attachment" json:"rq6_attachment"`
	Rq7Answer      bool           `gorm:"column:rq7_answer" json:"rq7_answer" form:"rq7_answer"`
	Rq7Attachments datatypes.JSON `gorm:"column:rq7_attachment" json:"rq7_attachment"`
	Cq1Answer      datatypes.JSON `gorm:"column:cq1_answer" json:"cq1_answer" form:"cq1_answer"`
	Cq1Attachments datatypes.JSON `gorm:"column:cq1_attachment" json:"cq1_attachment"`
	Cq2Answer      datatypes.JSON `gorm:"column:cq2_answer" json:"cq2_answer" form:"cq2_answer"`
	Cq2Attachments datatypes.JSON `gorm:"column:cq2_attachment" json:"cq2_attachment"`
	Cq3Answer      datatypes.JSON `gorm:"column:cq3_answer" json:"cq3_answer" form:"cq3_answer"`
	Cq3Attachments datatypes.JSON `gorm:"column:cq3_attachment" json:"cq3_attachment"`
	Cq4Answer      datatypes.JSON `gorm:"column:cq4_answer" json:"cq4_answer" form:"cq4_answer"`
	Cq4Attachments datatypes.JSON `gorm:"column:cq4_attachment" json:"cq4_attachment"`
	Cq5Answer      datatypes.JSON `gorm:"column:cq5_answer" json:"cq5_answer" form:"cq5_answer"`
	Cq5Attachments datatypes.JSON `gorm:"column:cq5_attachment" json:"cq5_attachment"`
	Cq6Answer      datatypes.JSON `gorm:"column:cq6_answer" json:"cq6_answer" form:"cq6_answer"`
	Cq6Attachments datatypes.JSON `gorm:"column:cq6_attachment" json:"cq6_attachment"`
	Cq7Answer      datatypes.JSON `gorm:"column:cq7_answer" json:"cq7_answer" form:"cq7_answer"`
	Cq7Attachments datatypes.JSON `gorm:"column:cq7_attachment" json:"cq7_attachment"`
	Cq8Answer      datatypes.JSON `gorm:"column:cq8_answer" json:"cq8_answer" form:"cq8_answer"`
	Cq8Attachments datatypes.JSON `gorm:"column:cq8_attachment" json:"cq8_attachment"`
	Cq9Answer      datatypes.JSON `gorm:"column:cq9_answer" json:"cq9_answer" form:"cq9_answer"`
	Cq9Attachments datatypes.JSON `gorm:"column:cq9_attachment" json:"cq9_attachment"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (AssessmentTrl) TableName() string {
	return "assessment_trls"
}
