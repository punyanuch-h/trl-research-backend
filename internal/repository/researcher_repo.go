package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

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
func (r *ResearcherRepo) Login(email string, password string) (*models.ResearcherInfo, error) {
	var researcher models.ResearcherInfo
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
func (r *ResearcherRepo) GetResearcherAll() ([]models.ResearcherInfo, error) {
	var researchers []models.ResearcherInfo
	err := r.DB.Find(&researchers).Error
	return researchers, err
}

// 🟢 GetResearcherByID - fetch one researcher by ID
func (r *ResearcherRepo) GetResearcherByID(researcherID string) (*models.ResearcherInfo, error) {
	var researcher models.ResearcherInfo
	err := r.DB.Where("researcher_id = ?", researcherID).First(&researcher).Error
	if err != nil {
		return nil, err
	}
	return &researcher, nil
}

// 🟢 GetResearcherByCaseID
func (r *ResearcherRepo) GetResearcherByCaseID(caseID string) (*models.ResearcherInfo, error) {
	// CaseInfo has researcher_id. Researcher does NOT have CaseID.
	var c models.CaseInfo
	if err := r.DB.Where("case_id = ?", caseID).First(&c).Error; err != nil {
		return nil, err
	}

	return r.GetResearcherByID(c.ResearcherID)
}

// 🟢 CreateResearcher - auto-generate ResearcherID and create new record
func (r *ResearcherRepo) CreateResearcher(researcher *models.ResearcherInfo) error {
	// find last ID to generate next
	var lastResearcher models.ResearcherInfo
	nextID := "RS-00001"
	err := r.DB.Order("researcher_id desc").First(&lastResearcher).Error
	if err == nil {
		lastID := lastResearcher.ResearcherID
		numStr := strings.TrimPrefix(lastID, "RS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("RS-%05d", n+1)
		}
	}

	researcher.ResearcherID = nextID
	now := time.Now()
	researcher.CreatedAt = now
	researcher.UpdatedAt = now

	return r.DB.Create(researcher).Error
}

// 🟢 UpdateResearcherByID - update with UpdatedAt
func (r *ResearcherRepo) UpdateResearcherByID(researcherID string, data *models.ResearcherInfo) error {
	data.UpdatedAt = time.Now()
	return r.DB.Model(&models.ResearcherInfo{}).Where("researcher_id = ?", researcherID).Updates(data).Error
}
