package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"gorm.io/gorm"
)

type AssessmentTrlRepo struct {
	DB *gorm.DB
}

func NewAssessmentTrlRepo(db *gorm.DB) AssessmentTrlRepository {
	return &AssessmentTrlRepo{DB: db}
}

// 🟢 GetAssessmentTrlAll
func (r *AssessmentTrlRepo) GetAssessmentTrlAll() ([]models.AssessmentTrl, error) {
	var assessments []models.AssessmentTrl
	err := r.DB.Find(&assessments).Error
	return assessments, err
}

// 🟢 GetAssessmentTrlByID
func (r *AssessmentTrlRepo) GetAssessmentTrlByID(id string) (*models.AssessmentTrl, error) {
	var a models.AssessmentTrl
	err := r.DB.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// 🟢 GetAssessmentTrlByCaseID
func (r *AssessmentTrlRepo) GetAssessmentTrlByCaseID(caseID string) (*models.AssessmentTrl, error) {
	var a models.AssessmentTrl
	err := r.DB.Where("case_id = ?", caseID).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("assessment trl with case_id %s not found", caseID)
		}
		return nil, err
	}
	return &a, nil
}

// 🟢 CreateAssessmentTrl - auto generate ID AS-<UUID>
func (r *AssessmentTrlRepo) CreateAssessmentTrl(a *models.AssessmentTrl) error {
	a.ID = utils.GenerateID("AS")
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	return r.DB.Create(a).Error
}

// 🟢 UpdateAssessmentTrlByID
func (r *AssessmentTrlRepo) UpdateAssessmentTrlByID(id string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return r.DB.Model(&models.AssessmentTrl{}).Where("id = ?", id).Updates(data).Error
}
