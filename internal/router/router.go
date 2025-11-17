package router

import (
	"net/http"
	"time"

	auth "trl-research-backend/internal/auth"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/handlers"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(gcsClient *storage.GCSClient) *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // ปิด debug log ของ Gin
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// ✅ CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://punyanuch-h.github.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Firestore repositories
	adminRepo := repository.NewAdminRepo(database.FirestoreClient)
	researcherRepo := repository.NewResearcherRepo(database.FirestoreClient)
	coordinatorRepo := repository.NewCoordinatorRepo(database.FirestoreClient)
	supporterRepo := repository.NewSupporterRepo(database.FirestoreClient)
	appointmentRepo := repository.NewAppointmentRepo(database.FirestoreClient)
	caseRepo := repository.NewCaseRepo(database.FirestoreClient)
	ipRepo := repository.NewIntellectualPropertyRepo(database.FirestoreClient)
	assessmentTrlRepo := repository.NewAssessmentTrlRepo(database.FirestoreClient)
	fileRepo := repository.NewFileRepo(database.FirestoreClient)

	// ✅ Handlers
	adminHandler := &handlers.AdminHandler{Repo: adminRepo}
	researcherHandler := &handlers.ResearcherHandler{Repo: researcherRepo}
	coordinatorHandler := &handlers.CoordinatorHandler{Repo: coordinatorRepo}
	supporterHandler := &handlers.SupporterHandler{Repo: supporterRepo}
	appointmentHandler := &handlers.AppointmentHandler{Repo: appointmentRepo}
	caseHandler := &handlers.CaseHandler{Repo: caseRepo}
	ipHandler := &handlers.IntellectualPropertyHandler{Repo: ipRepo}
	assessmentTrlHandler := &handlers.AssessmentTrlHandler{Repo: assessmentTrlRepo}
	presignHandler := &handlers.PresignHandler{GCS: gcsClient}
	fileHandler := &handlers.FileHandler{Repo: fileRepo}
	fileDownloadHandler := &handlers.FileDownloadHandler{FileRepo: fileRepo, GCS: gcsClient}

	// ✅ Auth Handlers
	loginHandler := &auth.LoginHandler{
		AdminRepo:      adminRepo,
		ResearcherRepo: researcherRepo,
	}
	forgotHandler := &auth.ForgotHandler{AdminRepo: *adminRepo}
	resetHandler := &auth.ResetHandler{AdminRepo: *adminRepo}

	// ✅ Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	// ✅ Public Auth
	r.POST("/auth/login", loginHandler.Login)
	r.POST("/auth/forgot-password", forgotHandler.ForgotPassword)
	r.POST("/auth/reset-password", resetHandler.ResetPassword)
	r.POST("/admin", adminHandler.CreateAdmin)

	// ✅ Protected APIs
	api := r.Group("/trl")
	api.Use(auth.AuthMiddleware())
	{
		api.GET("/admins", adminHandler.GetAllAdmins)
		api.GET("/admin/:id", adminHandler.GetAdminByID)
		api.GET("/admin/profile", adminHandler.GetAdminProfile)
		api.PATCH("/admin/:id", adminHandler.UpdateAdminProfileByID)

		api.GET("/researchers", researcherHandler.GetResearcherAll)
		api.GET("/researcher/:id", researcherHandler.GetResearcherByID)
		api.GET("/researcher/case/:id", researcherHandler.GetResearcherByCaseID)
		api.POST("/researcher", researcherHandler.CreateResearcher)
		api.PATCH("/researcher/:id", researcherHandler.UpdateResearcherProfileByID)
		api.GET("/researcher/profile", researcherHandler.GetResearcherProfile)

		api.GET("/coordinators", coordinatorHandler.GetCoordinatorAll)
		api.GET("/coordinator/:id", coordinatorHandler.GetCoordinatorByEmail)
		api.GET("/coordinator/case/:id", coordinatorHandler.GetCoordinatorByCaseID)
		api.POST("/coordinator", coordinatorHandler.CreateCoordinator)
		api.PATCH("/coordinator/:id", coordinatorHandler.UpdateCoordinatorByEmail)

		api.GET("/supporters", supporterHandler.GetSupporterAll)
		api.GET("/supporter/:id", supporterHandler.GetSupporterByID)
		api.GET("/supporter/case/:id", supporterHandler.GetSupporterByCaseID)
		api.POST("/supporter", supporterHandler.CreateSupporter)
		api.PATCH("/supporter/:id", supporterHandler.UpdateSupporterByID)

		api.GET("/appointments", appointmentHandler.GetAppointmentAll)
		api.GET("/appointment/:id", appointmentHandler.GetAppointmentByID)
		api.GET("/appointment/case/:id", appointmentHandler.GetAppointmentByCaseID)
		api.POST("/appointment", appointmentHandler.CreateAppointment)
		api.PATCH("/appointment/:id", appointmentHandler.UpdateAppointmentByID)

		api.GET("/cases", caseHandler.GetCaseAll)
		api.GET("/case/researcher/:id", caseHandler.GetCaseAllByResearcher_id)
		api.GET("/case/:id", caseHandler.GetCaseByID)
		api.POST("/case", caseHandler.CreateCase)
		api.PATCH("/case/:id", caseHandler.UpdateCaseByID)
		api.PATCH("/case/update-status/:id", caseHandler.UpdateCaseStatusByID)

		api.GET("/ips", ipHandler.GetIPAll)
		api.GET("/ip/:id", ipHandler.GetIPByID)
		api.GET("/ip/case/:id", ipHandler.GetIPByCaseID)
		api.POST("/ip", ipHandler.CreateIP)
		api.PATCH("/ip/:id", ipHandler.UpdateIPByID)

		api.GET("/assessment_trl", assessmentTrlHandler.GetAssessmentTrlAll)
		api.GET("/assessment_trl/:id", assessmentTrlHandler.GetAssessmentTrlByID)
		api.GET("/assessment_trl/case/:id", assessmentTrlHandler.GetAssessmentTrlByCaseID)
		api.POST("/assessment_trl", assessmentTrlHandler.CreateAssessmentTrl)
		api.PATCH("/assessment_trl/:id", assessmentTrlHandler.UpdateAssessmentTrlByID)

		// 🟢 File Management
		api.POST("/presign/upload", presignHandler.PresignUpload)
		api.POST("/file/upload", fileHandler.FileUploaded)
		api.GET("/file/download-url/:fileID", fileDownloadHandler.GetDownloadURL)
	}

	return r
}
