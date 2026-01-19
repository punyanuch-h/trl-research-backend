package models

import (
	"time"
	"trl-research-backend/internal/entity"
)

type ResearcherInfo struct {
	ResearcherID     string    `gorm:"primaryKey;column:researcher_id" json:"researcher_id"`
	AdminID          string    `gorm:"column:admin_id" json:"admin_id"`
	Prefix           string    `gorm:"column:prefix" json:"researcher_prefix"`
	AcademicPosition string    `gorm:"column:academic_position" json:"researcher_academic_position"`
	FirstName        string    `gorm:"column:first_name" json:"researcher_first_name"`
	LastName         string    `gorm:"column:last_name" json:"researcher_last_name"`
	Department       string    `gorm:"column:department" json:"researcher_department"`
	PhoneNumber      string    `gorm:"column:phone_number" json:"researcher_phone_number"`
	Email            string    `gorm:"column:email;unique" json:"researcher_email"`
	Password         string    `gorm:"column:password" json:"researcher_password"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ResearcherInfo) TableName() string {
	return "researchers"
}

func (r *ResearcherInfo) ToResponse() entity.ResearcherResponse {
	return entity.ResearcherResponse{
		ID:               r.ResearcherID,
		Prefix:           r.Prefix,
		AcademicPosition: r.AcademicPosition,
		FirstName:        r.FirstName,
		LastName:         r.LastName,
		Department:       r.Department,
		PhoneNumber:      r.PhoneNumber,
		Email:            r.Email,
	}
}
