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
func (r *CaseRepo) GetCaseAll() ([]models.Cases, error) {
	var cases []models.Cases
	err := r.DB.Find(&cases).Error
	return cases, err
}

// 🟢 GetCaseAllByResearcherID - fetch all cases for a researcher
func (r *CaseRepo) GetCaseAllByResearcherID(researcherID string) ([]models.Cases, error) {
	var cases []models.Cases
	err := r.DB.Where("researcher_id = ?", researcherID).Find(&cases).Error
	return cases, err
}

// 🟢 GetCaseByID
func (r *CaseRepo) GetCaseByID(caseID string) (*models.Cases, error) {
	var cs models.Cases
	err := r.DB.Where("id = ?", caseID).First(&cs).Error
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// 🟢 CreateCase - auto generate CaseID (CS-00001)
func (r *CaseRepo) CreateCase(cs *models.Cases) error {
	id, err := utils.GenerateID(r.DB, "cases", "CS")
	if err != nil {
		return err
	}
	cs.ID = id
	now := time.Now()
	cs.CreatedAt = now
	cs.UpdatedAt = now

	return r.DB.Create(cs).Error
}

// 🟢 UpdateCaseByID
func (r *CaseRepo) UpdateCaseByID(caseID string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Cases{}).Where("id = ?", caseID).Updates(data).Error
}

// 🟢 UpdateCaseStatusByID
func (r *CaseRepo) UpdateCaseStatusByID(caseID string, status bool) error {
	return r.DB.Model(&models.Cases{}).Where("id = ?", caseID).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}
