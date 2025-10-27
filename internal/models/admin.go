package models

import (
	"time"
	"trl-research-backend/internal/entity"
)

type AdminInfo struct {
	AdminID               string    `json:"admin_id" firestore:"admin_id"`
	AdminPrefix           string    `json:"admin_prefix" firestore:"admin_prefix"`
	AdminAcademicPosition string    `json:"admin_academic_position" firestore:"admin_academic_position"`
	AdminFirstName        string    `json:"admin_first_name" firestore:"admin_first_name"`
	AdminLastName         string    `json:"admin_last_name" firestore:"admin_last_name"`
	AdminDepartment       string    `json:"admin_department" firestore:"admin_department"`
	AdminPhoneNumber      string    `json:"admin_phone_number" firestore:"admin_phone_number"`
	AdminEmail            string    `json:"admin_email" firestore:"admin_email"`
	AdminPassword         string    `json:"admin_password" firestore:"admin_password"`
	CaseID                string    `json:"case_id" firestore:"case_id"`
	CreatedAt             time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" firestore:"updated_at"`
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
