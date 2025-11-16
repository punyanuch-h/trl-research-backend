package models

import "time"

type FileMetadata struct {
    ID              string    `json:"id" firestore:"id"`
    FileName        string    `json:"file_name" firestore:"file_name"`
    ObjectPath      string    `json:"object_path" firestore:"object_path"`
    Bucket          string    `json:"bucket" firestore:"bucket"`
    UploadedBy      string    `json:"uploaded_by" firestore:"uploaded_by"`
    UploadedAt      time.Time `json:"uploaded_at" firestore:"uploaded_at"`
    ContentType     string    `json:"content_type" firestore:"content_type"`
    BelongsToCaseID string    `json:"belongs_to_case_id" firestore:"belongs_to_case_id"` // optional
}
