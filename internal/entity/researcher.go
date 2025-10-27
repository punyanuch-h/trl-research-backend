package entity

type ResearcherResponse struct {
	ID               string `json:"id"`
	Prefix           string `json:"prefix"`
	AcademicPosition string `json:"academic_position"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Department       string `json:"department"`
	PhoneNumber      string `json:"phone_number"`
	Email            string `json:"email"`
}
