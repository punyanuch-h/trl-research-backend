package repository

import (
	"fmt"
	"time"

	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepo(db *gorm.DB) AdminRepository {
	return &AdminRepo{DB: db}
}

// 🟢 Get all admins
func (r *AdminRepo) GetAdminAll() ([]models.AdminInfo, error) {
	var admins []models.AdminInfo
	err := r.DB.Find(&admins).Error
	return admins, err
}

// 🟢 Get admin by ID
func (r *AdminRepo) GetAdminByID(adminID string) (*models.AdminInfo, error) {
	var admin models.AdminInfo
	err := r.DB.Where("admin_id = ?", adminID).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// 🟢 Get admin by email
func (r *AdminRepo) GetAdminByEmail(email string) (*models.AdminInfo, error) {
	var admin models.AdminInfo
	err := r.DB.Where("admin_email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// 🟢 Create admin (auto-generate AdminID)
func (r *AdminRepo) CreateAdmin(admin *models.AdminInfo) error {
	admin.AdminID = utils.GenerateID("AD")
	now := time.Now()
	admin.CreatedAt = now
	admin.UpdatedAt = now

	return r.DB.Create(admin).Error
}

// 🟢 Login with password verification
func (r *AdminRepo) Login(email string, password string) (*models.AdminInfo, error) {
	var admin models.AdminInfo
	err := r.DB.Where("admin_email = ?", email).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &admin, nil
}

// 🟢 Update password
func (r *AdminRepo) UpdatePasswordByEmail(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	return r.DB.Model(&models.AdminInfo{}).Where("admin_email = ?", email).Update("admin_password", string(hashedPassword)).Error
}

// 🟢 Update admin by ID
func (r *AdminRepo) UpdateAdminByID(adminID string, data *models.AdminInfo) error {
	data.UpdatedAt = time.Now()
	return r.DB.Model(&models.AdminInfo{}).Where("admin_id = ?", adminID).Updates(data).Error
}

// 🟢 Delete admin
func (r *AdminRepo) DeleteAdmin(email string) error {
	return r.DB.Where("admin_email = ?", email).Delete(&models.AdminInfo{}).Error
}
