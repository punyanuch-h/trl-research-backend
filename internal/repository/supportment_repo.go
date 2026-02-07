package repository

import (
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type SupportmentRepo struct {
	DB *gorm.DB
}

func NewSupportmentRepo(db *gorm.DB) SupportmentRepository {
	return &SupportmentRepo{DB: db}
}

// 🟢 GetSupportmentAll
func (r *SupportmentRepo) GetSupportmentAll() ([]models.Supportments, error) {
	var supportments []models.Supportments
	err := r.DB.Find(&supportments).Error
	return supportments, err
}

// 🟢 GetSupportmentByID
func (r *SupportmentRepo) GetSupportmentByID(supportmentID string) (*models.Supportments, error) {
	var supportment models.Supportments
	err := r.DB.Where("id = ?", supportmentID).First(&supportment).Error
	if err != nil {
		return nil, err
	}
	return &supportment, nil
}

// 🟢 GetSupportmentByCaseID
func (r *SupportmentRepo) GetSupportmentByCaseID(caseID string) (*models.Supportments, error) {
	var supportment models.Supportments
	err := r.DB.Where("case_id = ?", caseID).First(&supportment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &supportment, nil
}

// 🟢 CreateSupportment - auto-generate ID SP-00001
func (r *SupportmentRepo) CreateSupportment(supportment *models.Supportments) error {
	id, err := utils.GenerateID(r.DB, "supportments", "SP")
	if err != nil {
		return err
	}
	supportment.ID = id
	now := time.Now()
	supportment.CreatedAt = now
	supportment.UpdatedAt = now

	return r.DB.Create(supportment).Error
}

// 🟢 UpdateSupportmentByID
func (r *SupportmentRepo) UpdateSupportmentByID(supportmentID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Supportments{}).Where("id = ?", supportmentID).Updates(data).Error
}
