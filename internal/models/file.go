package models

import "time"

type FileMetadata struct {
	ID              string    `gorm:"primaryKey;column:id" json:"id"`
	FileName        string    `gorm:"column:file_name" json:"file_name"`
	ObjectPath      string    `gorm:"column:object_path" json:"object_path"`
	Bucket          string    `gorm:"column:bucket" json:"bucket"`
	UploadedBy      string    `gorm:"column:uploaded_by" json:"uploaded_by"`
	UploadedAt      time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	ContentType     string    `gorm:"column:content_type" json:"content_type"`
	BelongsToCaseID string    `gorm:"column:belongs_to_case_id" json:"belongs_to_case_id"` // optional
}

func (FileMetadata) TableName() string {
	return "file_metadatas"
}
