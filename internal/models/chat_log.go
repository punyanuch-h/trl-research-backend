package models

import (
	"time"

	"gorm.io/datatypes"
)

type ChatLogs struct {
	ID           string         `gorm:"primaryKey;column:id" json:"id"`
	AdminID      *string        `gorm:"column:admin_id;index" json:"admin_id,omitempty"`
	Admin        *Admins        `gorm:"foreignKey:AdminID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"admin,omitempty"`
	ResearcherID *string        `gorm:"column:researcher_id;index" json:"researcher_id,omitempty"`
	Researcher   *Researchers   `gorm:"foreignKey:ResearcherID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"researcher,omitempty"`
	History      datatypes.JSON `gorm:"column:history;type:text;not null" json:"history"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (ChatLogs) TableName() string {
	return "chat_logs"
}
