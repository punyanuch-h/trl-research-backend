package models

import (
	"time"
)

type RefreshToken struct {
	ID              string     `gorm:"primaryKey;column:id" json:"id"`
	UserID          string     `gorm:"column:user_id;not null" json:"user_id"`
	TokenHash       string     `gorm:"column:token_hash;not null;unique" json:"-"`
	ExpiryAt        time.Time  `gorm:"column:expiry_at;not null" json:"expiry_at"`
	RevokedAt       *time.Time `gorm:"column:revoked_at" json:"revoked_at"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UserType        string     `gorm:"column:user_type;not null" json:"user_type"` // admin or researcher
	ReplacedByToken string     `gorm:"column:replaced_by_token" json:"replaced_by_token"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
