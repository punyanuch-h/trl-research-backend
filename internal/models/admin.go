package models

import (
	"time"
	"trl-research-backend/internal/entity"
)

type AdminInfo struct {
	AdminID               string    `gorm:"primaryKey;column:admin_id" json:"admin_id"`
	AdminPrefix           string    `gorm:"column:admin_prefix" json:"admin_prefix"`
	AdminAcademicPosition string    `gorm:"column:admin_academic_position" json:"admin_academic_position"`
	AdminFirstName        string    `gorm:"column:admin_first_name" json:"admin_first_name"`
	AdminLastName         string    `gorm:"column:admin_last_name" json:"admin_last_name"`
	AdminDepartment       string    `gorm:"column:admin_department" json:"admin_department"`
	AdminPhoneNumber      string    `gorm:"column:admin_phone_number" json:"admin_phone_number"`
	AdminEmail            string    `gorm:"column:admin_email;unique" json:"admin_email"`
	AdminPassword         string    `gorm:"column:admin_password" json:"admin_password"`
	CaseID                string    `gorm:"column:case_id" json:"case_id"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AdminInfo) TableName() string {
	return "admins"
}

func (r *AdminInfo) ToResponse() entity.AdminResponse {
	return entity.AdminResponse{
		ID:               r.AdminID,
		Prefix:           r.AdminPrefix,
		AcademicPosition: r.AdminAcademicPosition,
		FirstName:        r.AdminFirstName,
		LastName:         r.AdminLastName,
		Department:       r.AdminDepartment,
		PhoneNumber:      r.AdminPhoneNumber,
		Email:            r.AdminEmail,
	}
}
