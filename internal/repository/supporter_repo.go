package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type SupporterRepo struct {
	DB *gorm.DB
}

func NewSupporterRepo(db *gorm.DB) SupporterRepository {
	return &SupporterRepo{DB: db}
}

// 🟢 GetSupporterAll
func (r *SupporterRepo) GetSupporterAll() ([]models.Supporter, error) {
	var supporters []models.Supporter
	err := r.DB.Find(&supporters).Error
	return supporters, err
}

// 🟢 GetSupporterByID
func (r *SupporterRepo) GetSupporterByID(supporterID string) (*models.Supporter, error) {
	var supporter models.Supporter
	err := r.DB.Where("supporter_id = ?", supporterID).First(&supporter).Error
	if err != nil {
		return nil, err
	}
	return &supporter, nil
}

// 🟢 GetSupporterByCaseID
func (r *SupporterRepo) GetSupporterByCaseID(caseID string) (*models.Supporter, error) {
	var supporter models.Supporter
	err := r.DB.Where("case_id = ?", caseID).First(&supporter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("supporter with case_id %s not found", caseID)
		}
		return nil, err
	}
	return &supporter, nil
}

// 🟢 CreateSupporter - auto-generate ID SP-00001
func (r *SupporterRepo) CreateSupporter(supporter *models.Supporter) error {
	var lastSupporter models.Supporter
	nextID := "SP-00001"
	err := r.DB.Order("supporter_id desc").First(&lastSupporter).Error
	if err == nil {
		lastID := lastSupporter.SupporterID
		numStr := strings.TrimPrefix(lastID, "SP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("SP-%05d", n+1)
		}
	}

	supporter.SupporterID = nextID
	now := time.Now()
	supporter.CreatedAt = now
	supporter.UpdatedAt = now

	return r.DB.Create(supporter).Error
}

// 🟢 UpdateSupporterByID
func (r *SupporterRepo) UpdateSupporterByID(supporterID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Supporter{}).Where("supporter_id = ?", supporterID).Updates(data).Error
}
