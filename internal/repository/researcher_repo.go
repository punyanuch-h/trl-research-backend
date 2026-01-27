package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ResearcherRepo struct {
	DB *gorm.DB
}

func NewResearcherRepo(db *gorm.DB) ResearcherRepository {
	return &ResearcherRepo{DB: db}
}

// 🟢 Login with password verification
func (r *ResearcherRepo) Login(email string, password string) (*models.Researchers, error) {
	var researcher models.Researchers
	err := r.DB.Where("email = ?", email).First(&researcher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("researcher not found")
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(researcher.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &researcher, nil
}

// 🟢 GetResearcherAll - fetch all researchers
func (r *ResearcherRepo) GetResearcherAll() ([]models.Researchers, error) {
	var researchers []models.Researchers
	err := r.DB.Find(&researchers).Error
	return researchers, err
}

// 🟢 GetResearcherByID - fetch one researcher by ID
func (r *ResearcherRepo) GetResearcherByID(researcherID string) (*models.Researchers, error) {
	var researcher models.Researchers
	err := r.DB.Where("id = ?", researcherID).First(&researcher).Error
	if err != nil {
		return nil, err
	}
	return &researcher, nil
}

// 🟢 GetResearcherByCaseID
func (r *ResearcherRepo) GetResearcherByCaseID(caseID string) (*models.Researchers, error) {
	// Cases has researcher_id. Researchers does NOT have CaseID.
	var c models.Cases
	if err := r.DB.Where("id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}

	return r.GetResearcherByID(c.ResearcherID)
}

// 🟢 CreateResearcher - auto-generate ResearcherID (RS-00001) and create new record
func (r *ResearcherRepo) CreateResearcher(researcher *models.Researchers) error {
	id, err := utils.GenerateID(r.DB, "researchers", "RS")
	if err != nil {
		return err
	}
	researcher.ID = id
	now := time.Now()
	researcher.CreatedAt = now
	researcher.UpdatedAt = now

	return r.DB.Create(researcher).Error
}

// 🟢 UpdateResearcherByID - update with UpdatedAt
func (r *ResearcherRepo) UpdateResearcherByID(researcherID string, data *models.Researchers) error {
	data.UpdatedAt = time.Now()
	return r.DB.Model(&models.Researchers{}).Where("id = ?", researcherID).Updates(data).Error
}
