package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type IntellectualPropertyRepo struct {
	DB *gorm.DB
}

func NewIntellectualPropertyRepo(db *gorm.DB) IntellectualPropertyRepository {
	return &IntellectualPropertyRepo{DB: db}
}

// 🟢 GetIPAll - fetch all intellectual property records
func (r *IntellectualPropertyRepo) GetIPAll() ([]models.IntellectualProperty, error) {
	var ips []models.IntellectualProperty
	err := r.DB.Find(&ips).Error
	return ips, err
}

// 🟢 GetIPByID
func (r *IntellectualPropertyRepo) GetIPByID(ipID string) (*models.IntellectualProperty, error) {
	var ip models.IntellectualProperty
	err := r.DB.Where("id = ?", ipID).First(&ip).Error
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// 🟢 GetIPByCaseID
func (r *IntellectualPropertyRepo) GetIPByCaseID(caseID string) (*models.IntellectualProperty, error) {
	var ip models.IntellectualProperty
	err := r.DB.Where("case_id = ?", caseID).First(&ip).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("intellectual property with case_id %s not found", caseID)
		}
		return nil, err
	}
	return &ip, nil
}

// 🟢 CreateIP - auto generate ID IP-00001
func (r *IntellectualPropertyRepo) CreateIP(ip *models.IntellectualProperty) error {
	var lastIP models.IntellectualProperty
	nextID := "IP-00001"
	err := r.DB.Order("id desc").First(&lastIP).Error
	if err == nil {
		lastID := lastIP.ID
		numStr := strings.TrimPrefix(lastID, "IP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("IP-%05d", n+1)
		}
	}

	ip.ID = nextID
	now := time.Now()
	ip.CreatedAt = now
	ip.UpdatedAt = now

	return r.DB.Create(ip).Error
}

// 🟢 UpdateIPByID
func (r *IntellectualPropertyRepo) UpdateIPByID(ipID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.IntellectualProperty{}).Where("id = ?", ipID).Updates(data).Error
}
