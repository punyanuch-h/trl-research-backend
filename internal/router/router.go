package router

import (
	"net/http"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/handlers"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 🔹 Initialize repository
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)
	researcherRepo := repository.NewResearcherRepo(database.FirestoreClient)
	coordinatorRepo := repository.NewCoordinatorRepo(database.FirestoreClient)
	supporterRepo := repository.NewSupporterRepo(database.FirestoreClient)

	appointmentRepo := repository.NewAppointmentRepo(database.FirestoreClient)
	caseRepo := repository.NewCaseRepo(database.FirestoreClient)
	ipRepo := repository.NewIntellectualPropertyRepo(database.FirestoreClient)
	assessmentTrlRepo := repository.NewAssessmentTrlRepo(database.FirestoreClient)

	// 🔹 Initialize handler
	adminHandler := &handlers.AdminHandler{Repo: adminRepo}
	researcherHandler := &handlers.ResearcherHandler{Repo: researcherRepo}
	coordinatorHandler := &handlers.CoordinatorHandler{Repo: coordinatorRepo}
	supporterHandler := &handlers.SupporterHandler{Repo: supporterRepo}

	appointmentHandler := &handlers.AppointmentHandler{Repo: appointmentRepo}
	caseHandler := &handlers.CaseHandler{Repo: caseRepo}
	ipHandler := &handlers.IntellectualPropertyHandler{Repo: ipRepo}
	assessmentTrlHandler := &handlers.AssessmentTrlHandler{Repo: assessmentTrlRepo}

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Admin
	r.GET("/admins", adminHandler.GetAllAdmins)
	r.GET("/admin/:id", adminHandler.GetAdminByID)
	r.POST("/admin", adminHandler.CreateAdmin)
	r.PATCH("/admin/:id", adminHandler.UpdateAdminByID)

	// Researcher
	r.GET("/researchers", researcherHandler.GetResearcherAll)
	r.GET("/researcher/:id", researcherHandler.GetResearcherByID)
	r.POST("/researcher", researcherHandler.CreateResearcher)
	r.PATCH("/researcher/:id", researcherHandler.UpdateResearcherByID)

	// Coordinator
	r.GET("/coordinators", coordinatorHandler.GetCoordinatorAll)
	r.GET("/coordinator/:id", coordinatorHandler.GetCoordinatorByEmail)
	r.POST("/coordinator", coordinatorHandler.CreateCoordinator)
	r.PATCH("/coordinator/:id", coordinatorHandler.UpdateCoordinatorByEmail)

	// Supporter
	r.GET("/supporters", supporterHandler.GetSupporterAll)
	r.GET("/supporter/:id", supporterHandler.GetSupporterByID)
	r.POST("/supporter", supporterHandler.CreateSupporter)
	r.PATCH("/supporter/:id", supporterHandler.UpdateSupporterByID)

	// Appointment
	r.GET("/appointments", appointmentHandler.GetAppointmentAll)
	r.GET("/appointment/:id", appointmentHandler.GetAppointmentByID)
	r.POST("/appointment", appointmentHandler.CreateAppointment)
	r.PATCH("/appintment/:id", appointmentHandler.UpdateAppointmentByID)

	// Case
	r.GET("/cases", caseHandler.GetCaseAll)
	r.GET("/case/:id", caseHandler.GetCaseByID)
	r.POST("/case", caseHandler.CreateCase)
	r.PATCH("/case/:id", caseHandler.UpdateCaseByID)

	// Intellectual Property
	r.GET("/ips", ipHandler.GetIPAll)
	r.GET("/ip/:id", ipHandler.GetIPByID)
	r.POST("/ip", ipHandler.CreateIP)
	r.PATCH("/ip/:id", ipHandler.UpdateIPByID)

	// Assessment TRL
	r.GET("/assessment_trl", assessmentTrlHandler.GetAssessmentTrlAll)
	r.GET("/assessment_trl/:id", assessmentTrlHandler.GetAssessmentTrlByID)
	r.POST("/assessment_trl", assessmentTrlHandler.CreateAssessmentTrl)
	r.PATCH("/assessment_trl/:id", assessmentTrlHandler.UpdateAssessmentTrlByID)

	return r
}
