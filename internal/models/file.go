package models

import "time"

type Files struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	CaseID      string    `gorm:"column:case_id;not null" json:"case_id"`
	Case        *Cases    `gorm:"foreignKey:CaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"case"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	ObjectPath  string    `gorm:"column:object_path;not null" json:"object_path"`
	Bucket      string    `gorm:"column:bucket" json:"bucket"`
	ContentType string    `gorm:"column:content_type" json:"content_type"`
	UploadedBy  string    `gorm:"column:uploaded_by;not null" json:"uploaded_by"`
	UploadedAt  time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
}

func (Files) TableName() string {
	return "files"
}
