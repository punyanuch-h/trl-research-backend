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
func (r *AdminRepo) GetAdminAll() ([]models.Admins, error) {
	var admins []models.Admins
	err := r.DB.Find(&admins).Error
	return admins, err
}

// 🟢 Get admin by ID
func (r *AdminRepo) GetAdminByID(adminID string) (*models.Admins, error) {
	var admin models.Admins
	err := r.DB.Where("id = ?", adminID).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// 🟢 Get admin by email
func (r *AdminRepo) GetAdminByEmail(email string) (*models.Admins, error) {
	var admin models.Admins
	err := r.DB.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// 🟢 Create admin (auto-generate AdminID AD-00001)
func (r *AdminRepo) CreateAdmin(admin *models.Admins) error {
	id, err := utils.GenerateID(r.DB, "admins", "AD")
	if err != nil {
		return err
	}
	admin.ID = id
	now := time.Now()
	admin.CreatedAt = now
	admin.UpdatedAt = now

	return r.DB.Create(admin).Error
}

// 🟢 Login with password verification
func (r *AdminRepo) Login(email string, password string) (*models.Admins, error) {
	var admin models.Admins
	err := r.DB.Where("email = ?", email).First(&admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
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
	return r.DB.Model(&models.Admins{}).Where("email = ?", email).Update("password", string(hashedPassword)).Error
}

// 🟢 Update admin by ID
func (r *AdminRepo) UpdateAdminByID(adminID string, data *models.Admins) error {
	data.UpdatedAt = time.Now()
	return r.DB.Model(&models.Admins{}).Where("id = ?", adminID).Updates(data).Error
}

// 🟢 Delete admin
func (r *AdminRepo) DeleteAdmin(email string) error {
	return r.DB.Where("email = ?", email).Delete(&models.Admins{}).Error
}
