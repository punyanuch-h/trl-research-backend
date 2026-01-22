package repository

import (
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

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

// 🟢 CreateCase - auto generate CaseID (CS-<UUID>)
func (r *CaseRepo) CreateCase(cs *models.CaseInfo) error {
	cs.CaseID = utils.GenerateID("CS")
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
func (r *CaseRepo) UpdateCaseStatusByID(caseID string, status bool) error {
	return r.DB.Model(&models.CaseInfo{}).Where("case_id = ?", caseID).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}
