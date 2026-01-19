package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type CoordinatorRepo struct {
	DB *gorm.DB
}

func NewCoordinatorRepo(db *gorm.DB) CoordinatorRepository {
	return &CoordinatorRepo{DB: db}
}

// 🟢 GetCoordinatorAll - fetch all coordinators
func (r *CoordinatorRepo) GetCoordinatorAll() ([]models.CoordinatorInfo, error) {
	var coordinators []models.CoordinatorInfo
	err := r.DB.Find(&coordinators).Error
	return coordinators, err
}

// 🟢 GetCoordinatorByEmail
func (r *CoordinatorRepo) GetCoordinatorByEmail(email string) (*models.CoordinatorInfo, error) {
	var coordinator models.CoordinatorInfo
	err := r.DB.Where("coordinator_email = ?", email).First(&coordinator).Error
	if err != nil {
		return nil, err
	}
	return &coordinator, nil
}

// 🟢 GetCoordinatorByCaseID
func (r *CoordinatorRepo) GetCoordinatorByCaseID(caseID string) (*models.CoordinatorInfo, error) {
	var coordinator models.CoordinatorInfo
	err := r.DB.Where("case_id = ?", caseID).First(&coordinator).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("coordinator with case_id %s not found", caseID)
		}
		return nil, err
	}
	return &coordinator, nil
}

// 🟢 CreateCoordinator
func (r *CoordinatorRepo) CreateCoordinator(coordinator *models.CoordinatorInfo) error {
	now := time.Now()
	coordinator.CreatedAt = now
	coordinator.UpdatedAt = now

	return r.DB.Create(coordinator).Error
}

// 🟢 UpdateCoordinatorByEmail
func (r *CoordinatorRepo) UpdateCoordinatorByEmail(email string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.CoordinatorInfo{}).Where("coordinator_email = ?", email).Updates(data).Error
}
