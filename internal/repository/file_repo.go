package repository

import (
	"context"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type FileRepo struct {
	DB *gorm.DB
}

func NewFileRepo(db *gorm.DB) FileRepository {
	return &FileRepo{DB: db}
}

func (r *FileRepo) SaveFile(ctx context.Context, file *models.Files) error {
	return r.DB.WithContext(ctx).Create(file).Error
}

func (r *FileRepo) GetFileByID(ctx context.Context, fileID string) (*models.Files, error) {
	var file models.Files
	err := r.DB.WithContext(ctx).Where("id = ?", fileID).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}
