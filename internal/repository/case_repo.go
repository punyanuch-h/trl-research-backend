package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"gorm.io/gorm"
)

type CaseRepo struct {
	DB *gorm.DB
}

func NewCaseRepo(db *gorm.DB) CaseRepository {
	return &CaseRepo{DB: db}
}

// 🟢 GetCaseAll - fetch all cases
func (r *CaseRepo) GetCaseAll() ([]models.CaseInfo, error) {
	var cases []models.CaseInfo
	err := r.DB.Find(&cases).Error
	return cases, err
}

// 🟢 GetCaseAllByResearcher_id - fetch all cases for a researcher
func (r *CaseRepo) GetCaseAllByResearcher_id(researcherID string) ([]models.CaseInfo, error) {
	var cases []models.CaseInfo
	err := r.DB.Where("researcher_id = ?", researcherID).Find(&cases).Error
	return cases, err
}

// 🟢 GetCaseByID
func (r *CaseRepo) GetCaseByID(caseID string) (*models.CaseInfo, error) {
	var cs models.CaseInfo
	err := r.DB.Where("case_id = ?", caseID).First(&cs).Error
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// 🟢 CreateCase - auto generate CaseID (CS-00001)
func (r *CaseRepo) CreateCase(cs *models.CaseInfo) error {
	var lastCase models.CaseInfo
	nextID := "CS-00001"
	err := r.DB.Order("case_id desc").First(&lastCase).Error
	if err == nil {
		lastID := lastCase.CaseID
		numStr := strings.TrimPrefix(lastID, "CS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("CS-%05d", n+1)
		}
	}

	cs.CaseID = nextID
	now := time.Now()
	cs.CreatedAt = now
	cs.UpdatedAt = now

	return r.DB.Create(cs).Error
}

// 🟢 UpdateCaseByID
func (r *CaseRepo) UpdateCaseByID(caseID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.CaseInfo{}).Where("case_id = ?", caseID).Updates(data).Error
}

// 🟢 UpdateCaseStatusByID
func (r *CaseRepo) UpdateCaseStatusByID(caseID string, status string) error {
	return r.DB.Model(&models.CaseInfo{}).Where("case_id = ?", caseID).Update("status", status).Error
}
