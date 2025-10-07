package router

import (
	"net/http"

	"trl-research-backend/internal/database"
	"trl-research-backend/internal/handlers"
	auth "trl-research-backend/internal/auth"
	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Repos
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)
	researcherRepo := repository.NewResearcherRepo(database.FirestoreClient)
	coordinatorRepo := repository.NewCoordinatorRepo(database.FirestoreClient)
	supporterRepo := repository.NewSupporterRepo(database.FirestoreClient)
	appointmentRepo := repository.NewAppointmentRepo(database.FirestoreClient)
	caseRepo := repository.NewCaseRepo(database.FirestoreClient)
	ipRepo := repository.NewIntellectualPropertyRepo(database.FirestoreClient)
	assessmentTrlRepo := repository.NewAssessmentTrlRepo(database.FirestoreClient)

	// CRUD Handlers
	adminHandler := &handlers.AdminHandler{Repo: adminRepo}
	researcherHandler := &handlers.ResearcherHandler{Repo: researcherRepo}
	coordinatorHandler := &handlers.CoordinatorHandler{Repo: coordinatorRepo}
	supporterHandler := &handlers.SupporterHandler{Repo: supporterRepo}
	appointmentHandler := &handlers.AppointmentHandler{Repo: appointmentRepo}
	caseHandler := &handlers.CaseHandler{Repo: caseRepo}
	ipHandler := &handlers.IntellectualPropertyHandler{Repo: ipRepo}
	assessmentTrlHandler := &handlers.AssessmentTrlHandler{Repo: assessmentTrlRepo}

	// Auth Handlers (Gin version)
	loginHandler := &auth.LoginHandler{AdminRepo: *adminRepo}
	forgotHandler := &auth.ForgotHandler{AdminRepo: *adminRepo}
	resetHandler := &auth.ResetHandler{AdminRepo: *adminRepo}

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// Public Auth
	r.POST("/auth/login", loginHandler.Login)
	r.POST("/auth/forgot-password", forgotHandler.ForgotPassword)
	r.POST("/auth/reset-password", resetHandler.ResetPassword)
	r.POST("/admin", adminHandler.CreateAdmin)

	// Protected
	api := r.Group("/trl")
	// api.Use(auth.AuthMiddleware())
	{
		// Admin
		api.GET("/admins", adminHandler.GetAllAdmins)
		api.GET("/admin/:id", adminHandler.GetAdminByID)

		// Researcher
		api.GET("/researchers", researcherHandler.GetResearcherAll)
		api.GET("/researcher/:id", researcherHandler.GetResearcherByID)
		api.POST("/researcher", researcherHandler.CreateResearcher)
		api.PATCH("/researcher/:id", researcherHandler.UpdateResearcherByID)

		// Coordinator
		api.GET("/coordinators", coordinatorHandler.GetCoordinatorAll)
		api.GET("/coordinator/:id", coordinatorHandler.GetCoordinatorByEmail)
		api.POST("/coordinator", coordinatorHandler.CreateCoordinator)
		api.PATCH("/coordinator/:id", coordinatorHandler.UpdateCoordinatorByEmail)

		// Supporter
		api.GET("/supporters", supporterHandler.GetSupporterAll)
		api.GET("/supporter/:id", supporterHandler.GetSupporterByID)
		api.POST("/supporter", supporterHandler.CreateSupporter)
		api.PATCH("/supporter/:id", supporterHandler.UpdateSupporterByID)

		// Appointment
		api.GET("/appointments", appointmentHandler.GetAppointmentAll)
		api.GET("/appointment/:id", appointmentHandler.GetAppointmentByID)
		api.POST("/appointment", appointmentHandler.CreateAppointment)
		api.PATCH("/appointment/:id", appointmentHandler.UpdateAppointmentByID)

		// Case
		api.GET("/cases", caseHandler.GetCaseAll)
		api.GET("/case/:id", caseHandler.GetCaseByID)
		api.POST("/case", caseHandler.CreateCase)
		api.PATCH("/case/:id", caseHandler.UpdateCaseByID)

		// IP
		api.GET("/ips", ipHandler.GetIPAll)
		api.GET("/ip/:id", ipHandler.GetIPByID)
		api.POST("/ip", ipHandler.CreateIP)
		api.PATCH("/ip/:id", ipHandler.UpdateIPByID)

		// Assessment TRL
		api.GET("/assessment_trl", assessmentTrlHandler.GetAssessmentTrlAll)
		api.GET("/assessment_trl/:id", assessmentTrlHandler.GetAssessmentTrlByID)
		api.POST("/assessment_trl", assessmentTrlHandler.CreateAssessmentTrl)
		api.PATCH("/assessment_trl/:id", assessmentTrlHandler.UpdateAssessmentTrlByID)
	}

	return r
}
