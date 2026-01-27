package models

import (
	"time"
	"trl-research-backend/internal/entity"
)

type Admins struct {
	ID               string    `gorm:"primaryKey;column:id" json:"id"`
	Prefix           string    `gorm:"column:prefix;not null" json:"prefix"`
	AcademicPosition string    `gorm:"column:academic_position" json:"academic_position"`
	FirstName        string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName         string    `gorm:"column:last_name;not null" json:"last_name"`
	Department       string    `gorm:"column:department;not null" json:"department"`
	PhoneNumber      string    `gorm:"column:phone_number;not null" json:"phone_number"`
	Email            string    `gorm:"column:email;not null;unique" json:"email"`
	Password         string    `gorm:"column:password;not null" json:"-"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Admins) TableName() string {
	return "admins"
}

func (r *Admins) ToResponse() entity.AdminResponse {
	return entity.AdminResponse{
		ID:               r.ID,
		Prefix:           r.Prefix,
		AcademicPosition: r.AcademicPosition,
		FirstName:        r.FirstName,
		LastName:         r.LastName,
		Department:       r.Department,
		PhoneNumber:      r.PhoneNumber,
		Email:            r.Email,
	}
}
