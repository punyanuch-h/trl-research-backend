package repository

import (
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type IntellectualPropertyRepo struct {
	DB *gorm.DB
}

func NewIntellectualPropertyRepo(db *gorm.DB) IntellectualPropertyRepository {
	return &IntellectualPropertyRepo{DB: db}
}

// 🟢 GetIPAll - fetch all intellectual property records
func (r *IntellectualPropertyRepo) GetIPAll() ([]models.IntellectualProperties, error) {
	var ips []models.IntellectualProperties
	err := r.DB.Find(&ips).Error
	return ips, err
}

// 🟢 GetIPByID
func (r *IntellectualPropertyRepo) GetIPByID(ipID string) (*models.IntellectualProperties, error) {
	var ip models.IntellectualProperties
	err := r.DB.Where("id = ?", ipID).First(&ip).Error
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// 🟢 GetIPByCaseID
func (r *IntellectualPropertyRepo) GetIPByCaseID(caseID string) ([]models.IntellectualProperties, error) {
	var ip []models.IntellectualProperties
	err := r.DB.Where("case_id = ?", caseID).Find(&ip).Error
	return ip, err
}

// 🟢 CreateIP - auto generate ID IP-00001
func (r *IntellectualPropertyRepo) CreateIP(ip *models.IntellectualProperties) error {
	id, err := utils.GenerateID(r.DB, "intellectual_properties", "IP")
	if err != nil {
		return err
	}
	ip.ID = id
	now := time.Now()
	ip.CreatedAt = now
	ip.UpdatedAt = now

	return r.DB.Create(ip).Error
}

// 🟢 UpdateIPByID
func (r *IntellectualPropertyRepo) UpdateIPByID(ipID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.IntellectualProperties{}).Where("id = ?", ipID).Updates(data).Error
}
