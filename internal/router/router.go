package router

import (
	"net/http"
	"time"

	auth "trl-research-backend/internal/auth"
	"trl-research-backend/internal/config"
	"trl-research-backend/internal/cron"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/handlers"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(gcsClient *storage.GCSClient, cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // ปิด debug log ของ Gin
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"0.0.0.0/0"})

	// ✅ CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://punyanuch-h.github.io", "https://trl-research-frontend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Postgres (GORM) repositories
	adminRepo := repository.NewAdminRepo(database.DB)
	researcherRepo := repository.NewResearcherRepo(database.DB)
	coordinatorRepo := repository.NewCoordinatorRepo(database.DB)
	supportmentRepo := repository.NewSupportmentRepo(database.DB)
	appointmentRepo := repository.NewAppointmentRepo(database.DB)
	caseRepo := repository.NewCaseRepo(database.DB)
	ipRepo := repository.NewIntellectualPropertyRepo(database.DB)
	assessmentRepo := repository.NewAssessmentRepo(database.DB)

	// ✅ Handlers
	adminHandler := &handlers.AdminHandler{Repo: adminRepo}
	researcherHandler := &handlers.ResearcherHandler{Repo: researcherRepo}
	coordinatorHandler := &handlers.CoordinatorHandler{Repo: coordinatorRepo}
	supportmentHandler := &handlers.SupportmentHandler{Repo: supportmentRepo}
	appointmentHandler := &handlers.AppointmentHandler{
		Repo: appointmentRepo,
		Cfg:  cfg,
	}
	caseHandler := &handlers.CaseHandler{
		Repo: caseRepo,
		GCS:  gcsClient,
	}
	ipHandler := &handlers.IntellectualPropertyHandler{Repo: ipRepo, GCS: gcsClient}
	assessmentHandler := &handlers.AssessmentHandler{Repo: assessmentRepo, GCS: gcsClient}
	presignHandler := &handlers.PresignHandler{GCS: gcsClient}
	fileDownloadHandler := &handlers.FileDownloadHandler{GCS: gcsClient}

	// ✅ Auth Handlers
	loginHandler := &auth.LoginHandler{
		AdminRepo:      adminRepo,
		ResearcherRepo: researcherRepo,
	}
	forgotHandler := &auth.ForgotHandler{
		AdminRepo:      adminRepo,
		ResearcherRepo: researcherRepo,
	}
	resetHandler := &auth.ResetHandler{
		AdminRepo:      adminRepo,
		ResearcherRepo: researcherRepo,
	}

	// ✅ Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// ✅ Public Auth
	r.POST("/auth/login", loginHandler.Login)
	r.POST("/auth/forget-password", forgotHandler.ForgetPassword)
	r.POST("/researcher", researcherHandler.CreateResearcher)

	// ✅ Protected APIs
	api := r.Group("/trl")
	api.Use(auth.AuthMiddleware())
	{
		api.POST("/auth/reset-password", resetHandler.ResetPassword)

		api.GET("/admins", auth.RequireRoles("admin"), adminHandler.GetAllAdmins)
		api.GET("/admin/:id", auth.RequireRoles("admin"), adminHandler.GetAdminByID)
		api.GET("/admin/profile", auth.RequireRoles("admin"), adminHandler.GetAdminProfile)
		api.POST("/admin", auth.RequireRoles("admin"), adminHandler.CreateAdmin)
		api.PATCH("/admin/:id", auth.RequireRoles("admin"), adminHandler.UpdateAdminProfileByID)

		api.GET("/researchers", researcherHandler.GetResearcherAll)
		api.GET("/researcher/:id", researcherHandler.GetResearcherByID)
		api.GET("/researcher/case/:id", researcherHandler.GetResearcherByCaseID)
		api.PATCH("/researcher/:id", researcherHandler.UpdateResearcherProfileByID)
		api.GET("/researcher/profile", researcherHandler.GetResearcherProfile)

		api.GET("/coordinators", coordinatorHandler.GetCoordinatorAll)
		api.GET("/coordinator/:id", coordinatorHandler.GetCoordinatorByEmail)
		api.GET("/coordinator/case/:id", coordinatorHandler.GetCoordinatorByCaseID)
		api.POST("/coordinator", coordinatorHandler.CreateCoordinator)
		api.PATCH("/coordinator/:id", auth.RequireRoles("admin"), coordinatorHandler.UpdateCoordinatorByEmail)

		api.GET("/supportments", supportmentHandler.GetSupportmentAll)
		api.GET("/supportment/:id", supportmentHandler.GetSupportmentByID)
		api.GET("/supportment/case/:id", supportmentHandler.GetSupportmentByCaseID)
		api.POST("/supportment", supportmentHandler.CreateSupportment)
		api.PATCH("/supportment/:id", auth.RequireRoles("admin"), supportmentHandler.UpdateSupportmentByID)

		api.GET("/appointments", appointmentHandler.GetAppointmentAll)
		api.GET("/appointment/:id", appointmentHandler.GetAppointmentByID)
		api.GET("/appointment/case/:id", appointmentHandler.GetAppointmentByCaseID)
		api.POST("/appointment", auth.RequireRoles("admin"), appointmentHandler.CreateAppointment)
		api.PATCH("/appointment/:id", auth.RequireRoles("admin"), appointmentHandler.UpdateAppointmentByID)

		// 🟢 Appointment Notifications
		api.GET("/notifications/appointments", appointmentHandler.GetNotifications)
		api.PATCH("/notifications/appointments/read-all", appointmentHandler.MarkAllAsRead)
		api.PATCH("/notifications/appointments/:id/read", appointmentHandler.MarkAsRead)

		api.GET("/cases", caseHandler.GetCaseAll)
		api.GET("/case/researcher/:id", caseHandler.GetCaseAllByResearcherID)
		api.GET("/case/:id", caseHandler.GetCaseByID)
		api.POST("/case", caseHandler.CreateCase)
		api.PATCH("/case/:id", auth.RequireRoles("admin"), caseHandler.UpdateCaseByID)
		api.PATCH("/case/update-status/:id", auth.RequireRoles("admin"), caseHandler.UpdateCaseStatusByID)

		api.GET("/ips", ipHandler.GetIPAll)
		api.GET("/ip/:id", ipHandler.GetIPByID)
		api.GET("/ip/case/:id", ipHandler.GetIPByCaseID)
		api.POST("/ip", ipHandler.CreateIP)
		api.PATCH("/ip/:id", auth.RequireRoles("admin"), ipHandler.UpdateIPByID)

		api.GET("/assessments", assessmentHandler.GetAssessmentAll)
		api.GET("/assessment/:id", assessmentHandler.GetAssessmentByID)
		api.GET("/assessment/case/:id", assessmentHandler.GetAssessmentByCaseID)
		api.POST("/assessment", assessmentHandler.CreateAssessment)
		api.PATCH("/assessment/:id", auth.RequireRoles("admin"), assessmentHandler.UpdateAssessmentByID)
		// 🟢 File Management
		api.POST("/presign/upload", presignHandler.PresignUpload)
		api.GET("/file/download", fileDownloadHandler.GetDownloadURL)
	}

	// ✅ Start Cron Jobs
	reminderCron := cron.NewReminderCron(appointmentRepo, adminRepo, cfg)
	reminderCron.Start()

	return r
}
