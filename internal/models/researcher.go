package models

import (
	"time"
)

type ResearcherInfo struct {
	ResearcherID               string    `json:"researcher_id" firestore:"researcher_id"`
	AdminID                    string    `json:"admin_id" firestore:"admin_id"`
	ResearcherPrefix           string    `json:"researcher_prefix" firestore:"researcher_prefix"`
	ResearcherAcademicPosition string    `json:"researcher_academic_position" firestore:"researcher_academic_position"`
	ResearcherFirstName        string    `json:"researcher_first_name" firestore:"researcher_first_name"`
	ResearcherLastName         string    `json:"researcher_last_name" firestore:"researcher_last_name"`
	ResearcherDepartment       string    `json:"researcher_department" firestore:"researcher_department"`
	ResearcherPhoneNumber      string    `json:"researcher_phone_number" firestore:"researcher_phone_number"`
	ResearcherEmail            string    `json:"researcher_email" firestore:"researcher_email"`
	CreatedAt                  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" firestore:"updated_at"`
}
