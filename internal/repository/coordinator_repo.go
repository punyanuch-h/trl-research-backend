package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type CoordinatorRepo struct {
	DB *gorm.DB
}

func NewCoordinatorRepo(db *gorm.DB) CoordinatorRepository {
	return &CoordinatorRepo{DB: db}
}

// 🟢 GetCoordinatorAll - fetch all coordinators
func (r *CoordinatorRepo) GetCoordinatorAll() ([]models.Coordinators, error) {
	var coordinators []models.Coordinators
	err := r.DB.Find(&coordinators).Error
	return coordinators, err
}

// 🟢 GetCoordinatorByEmail
func (r *CoordinatorRepo) GetCoordinatorByEmail(email string) (*models.Coordinators, error) {
	var coordinator models.Coordinators
	err := r.DB.Where("email = ?", email).First(&coordinator).Error
	if err != nil {
		return nil, err
	}
	return &coordinator, nil
}

// 🟢 GetCoordinatorByCaseID
func (r *CoordinatorRepo) GetCoordinatorByCaseID(caseID string) (*models.Coordinators, error) {
    var caseData models.Cases
    if err := r.DB.Where("id = ?", caseID).First(&caseData).Error; err != nil {
        return nil, fmt.Errorf("case not found: %v", err)
    }
    var coordinator models.Coordinators
    if err := r.DB.Where("email = ?", caseData.CoordinatorEmail).First(&coordinator).Error; err != nil {
        return nil, fmt.Errorf("coordinator not found: %v", err)
    }
    return &coordinator, nil
}

// 🟢 CreateCoordinator - auto generate CoordinatorID (C-<UUID>) or update if email exists
func (r *CoordinatorRepo) CreateCoordinator(coordinator *models.Coordinators) error {
	// Check if a coordinator with this email already exists
	var existing models.Coordinators
	if err := r.DB.Where("email = ?", coordinator.Email).First(&existing).Error; err == nil {
		// If exists, update the existing record and reuse its ID
		coordinator.ID = existing.ID
		now := time.Now()
		coordinator.UpdatedAt = now
		return r.DB.Model(&existing).Updates(coordinator).Error
	}

	coordinator.ID = utils.GenerateID("CO")
	now := time.Now()
	coordinator.CreatedAt = now
	coordinator.UpdatedAt = now

	return r.DB.Create(coordinator).Error
}

// 🟢 UpdateCoordinatorByEmail
func (r *CoordinatorRepo) UpdateCoordinatorByEmail(email string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Coordinators{}).Where("email = ?", email).Updates(data).Error
}
