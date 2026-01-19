package repository

import (
	"context"
	"trl-research-backend/internal/models"
)

type AdminRepository interface {
	GetAdminAll() ([]models.AdminInfo, error)
	GetAdminByID(adminID string) (*models.AdminInfo, error)
	GetAdminByEmail(email string) (*models.AdminInfo, error)
	CreateAdmin(admin *models.AdminInfo) error
	Login(email string, password string) (*models.AdminInfo, error)
	UpdatePasswordByEmail(email string, password string) error
	UpdateAdminByID(adminID string, data *models.AdminInfo) error
	DeleteAdmin(email string) error
}

type ResearcherRepository interface {
	Login(email string, password string) (*models.ResearcherInfo, error)
	GetResearcherAll() ([]models.ResearcherInfo, error)
	GetResearcherByID(researcherID string) (*models.ResearcherInfo, error)
	GetResearcherByCaseID(caseID string) (*models.ResearcherInfo, error)
	CreateResearcher(researcher *models.ResearcherInfo) error
	UpdateResearcherByID(researcherID string, data *models.ResearcherInfo) error
}

type CaseRepository interface {
	GetCaseAll() ([]models.CaseInfo, error)
	GetCaseAllByResearcher_id(researcherID string) ([]models.CaseInfo, error)
	GetCaseByID(caseID string) (*models.CaseInfo, error)
	CreateCase(cs *models.CaseInfo) error
	UpdateCaseByID(caseID string, data map[string]interface{}) error
	UpdateCaseStatusByID(caseID string, status string) error
}

type SupporterRepository interface {
	GetSupporterAll() ([]models.Supporter, error)
	GetSupporterByID(supporterID string) (*models.Supporter, error)
	GetSupporterByCaseID(caseID string) (*models.Supporter, error)
	CreateSupporter(supporter *models.Supporter) error
	UpdateSupporterByID(supporterID string, data map[string]interface{}) error
}

type AppointmentRepository interface {
	GetAppointmentAll() ([]models.Appointment, error)
	GetAppointmentByID(appointmentID string) (*models.Appointment, error)
	GetAppointmentByCaseID(caseID string) ([]models.Appointment, error)
	CreateAppointment(ap *models.Appointment) error
	UpdateAppointmentByID(appointmentID string, data map[string]interface{}) error
}

type CoordinatorRepository interface {
	GetCoordinatorAll() ([]models.CoordinatorInfo, error)
	GetCoordinatorByEmail(email string) (*models.CoordinatorInfo, error)
	GetCoordinatorByCaseID(caseID string) (*models.CoordinatorInfo, error)
	CreateCoordinator(coordinator *models.CoordinatorInfo) error
	UpdateCoordinatorByEmail(email string, data map[string]interface{}) error
}

type AssessmentTrlRepository interface {
	GetAssessmentTrlAll() ([]models.AssessmentTrl, error)
	GetAssessmentTrlByID(id string) (*models.AssessmentTrl, error)
	GetAssessmentTrlByCaseID(caseID string) (*models.AssessmentTrl, error)
	CreateAssessmentTrl(a *models.AssessmentTrl) error
	UpdateAssessmentTrlByID(id string, data map[string]interface{}) error
}

type IntellectualPropertyRepository interface {
	GetIPAll() ([]models.IntellectualProperty, error)
	GetIPByID(ipID string) (*models.IntellectualProperty, error)
	GetIPByCaseID(caseID string) (*models.IntellectualProperty, error)
	CreateIP(ip *models.IntellectualProperty) error
	UpdateIPByID(ipID string, data map[string]interface{}) error
}

type FileRepository interface {
	SaveFile(ctx context.Context, file *models.FileMetadata) error
	GetFileByID(ctx context.Context, fileID string) (*models.FileMetadata, error)
}
