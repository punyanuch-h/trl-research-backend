package repository

import (
    "context"
    "trl-research-backend/internal/models"

    "cloud.google.com/go/firestore"
)

type FileRepo struct {
    client *firestore.Client
}

func NewFileRepo(client *firestore.Client) *FileRepo {
    return &FileRepo{client}
}

func (r *FileRepo) SaveFile(ctx context.Context, file *models.FileMetadata) error {
    _, err := r.client.
        Collection("files").
        Doc(file.ID).
        Set(ctx, file)
    return err
}

func (r *FileRepo) GetFileByID(ctx context.Context, fileID string) (*models.FileMetadata, error) {
	doc, err := r.client.Collection("files").Doc(fileID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var file models.FileMetadata
	if err := doc.DataTo(&file); err != nil {
		return nil, err
	}

	return &file, nil
}
