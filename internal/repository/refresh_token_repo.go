package repository

import (
	"time"
	"trl-research-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepo struct {
	DB *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenRepo{DB: db}
}

func (r *RefreshTokenRepo) Create(token *models.RefreshToken) error {
	if token.ID == "" {
		token.ID = uuid.New().String()
	}
	token.CreatedAt = time.Now()
	return r.DB.Create(token).Error
}

func (r *RefreshTokenRepo) GetByHash(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.DB.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *RefreshTokenRepo) Revoke(id string) error {
	now := time.Now()
	return r.DB.Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepo) RevokeAllForUser(userID string) error {
	now := time.Now()
	return r.DB.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepo) Update(token *models.RefreshToken) error {
	return r.DB.Save(token).Error
}

func (r *RefreshTokenRepo) DeleteExpired() error {
	now := time.Now()
	// Also delete revoked tokens that are older than some threshold? 
	// For now, simple delete expired.
	return r.DB.Where("expiry_at < ?", now).Delete(&models.RefreshToken{}).Error
}
