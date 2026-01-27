package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type AssessmentRepo struct {
	DB *gorm.DB
}

func NewAssessmentRepo(db *gorm.DB) AssessmentRepository {
	return &AssessmentRepo{DB: db}
}

// 🟢 GetAssessmentAll
func (r *AssessmentRepo) GetAssessmentAll() ([]models.Assessments, error) {
	var assessments []models.Assessments
	err := r.DB.Find(&assessments).Error
	return assessments, err
}

// 🟢 GetAssessmentByID
func (r *AssessmentRepo) GetAssessmentByID(id string) (*models.Assessments, error) {
	var a models.Assessments
	err := r.DB.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// 🟢 GetAssessmentByCaseID
func (r *AssessmentRepo) GetAssessmentByCaseID(caseID string) (*models.Assessments, error) {
	var a models.Assessments
	err := r.DB.Where("case_id = ?", caseID).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment with case_id %s not found", caseID)
		}
		return nil, err
	}
	return &a, nil
}

// 🟢 CreateAssessment - auto generate ID AS-<UUID>
func (r *AssessmentRepo) CreateAssessment(a *models.Assessments) error {
	a.ID = utils.GenerateID("AS")
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	return r.DB.Create(a).Error
}

// 🟢 UpdateAssessmentByID
func (r *AssessmentRepo) UpdateAssessmentByID(id string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.Assessments{}).Where("id = ?", id).Updates(data).Error
}
