package repository

import (
	"context"
	"trl-research-backend/internal/models"
)

type AdminRepository interface {
	GetAdminAll() ([]models.Admins, error)
	GetAdminByID(adminID string) (*models.Admins, error)
	GetAdminByEmail(email string) (*models.Admins, error)
	CreateAdmin(admin *models.Admins) error
	Login(email string, password string) (*models.Admins, error)
	UpdatePasswordByEmail(email string, password string) error
	UpdateAdminByID(adminID string, data *models.Admins) error
	DeleteAdmin(email string) error
}

type ResearcherRepository interface {
	Login(email string, password string) (*models.Researchers, error)
	UpdatePasswordByEmail(email string, password string) error
	GetResearcherAll() ([]models.Researchers, error)
	GetResearcherByID(researcherID string) (*models.Researchers, error)
	GetResearcherByCaseID(caseID string) (*models.Researchers, error)
	CreateResearcher(researcher *models.Researchers) error
	UpdateResearcherByID(researcherID string, data *models.Researchers) error
}

type CaseRepository interface {
	GetCaseAll() ([]models.Cases, error)
	GetCaseAllByResearcherID(researcherID string) ([]models.Cases, error)
	GetCaseByID(caseID string) (*models.Cases, error)
	CreateCase(cs *models.Cases) error
	UpdateCaseByID(caseID string, data map[string]interface{}) error
	UpdateCaseStatusByID(caseID string, status bool) error
}

type SupportmentRepository interface {
	GetSupportmentAll() ([]models.Supportments, error)
	GetSupportmentByID(supportmentID string) (*models.Supportments, error)
	GetSupportmentByCaseID(caseID string) (*models.Supportments, error)
	CreateSupportment(supportment *models.Supportments) error
	UpdateSupportmentByID(supportmentID string, data map[string]interface{}) error
}

type AppointmentRepository interface {
	GetAppointmentAll() ([]models.Appointments, error)
	GetAppointmentByID(appointmentID string) (*models.Appointments, error)
	GetAppointmentWithDetails(appointmentID string) (*models.Appointments, error)
	GetAppointmentByCaseID(caseID string) ([]models.Appointments, error)
	CreateAppointment(ap *models.Appointments) error
	UpdateAppointmentByID(appointmentID string, data map[string]interface{}) error
}

type CoordinatorRepository interface {
	GetCoordinatorAll() ([]models.Coordinators, error)
	GetCoordinatorByEmail(email string) (*models.Coordinators, error)
	GetCoordinatorByCaseID(caseID string) (*models.Coordinators, error)
	CreateCoordinator(coordinator *models.Coordinators) error
	UpdateCoordinatorByEmail(email string, data map[string]interface{}) error
}

type AssessmentRepository interface {
	GetAssessmentAll() ([]models.Assessments, error)
	GetAssessmentByID(id string) (*models.Assessments, error)
	GetAssessmentByCaseID(caseID string) (*models.Assessments, error)
	CreateAssessment(a *models.Assessments) error
	UpdateAssessmentByID(id string, data map[string]interface{}) error
}

type IntellectualPropertyRepository interface {
	GetIPAll() ([]models.IntellectualProperties, error)
	GetIPByID(ipID string) (*models.IntellectualProperties, error)
	GetIPByCaseID(caseID string) ([]models.IntellectualProperties, error)
	CreateIP(ip *models.IntellectualProperties) error
	UpdateIPByID(ipID string, data map[string]interface{}) error
}

type FileRepository interface {
	SaveFile(ctx context.Context, file *models.Files) error
	GetFileByID(ctx context.Context, fileID string) (*models.Files, error)
}
