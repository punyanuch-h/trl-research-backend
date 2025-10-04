package router

import (
	"github.com/gin-gonic/gin"
	"trl-research-backend/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ---------------- Root test endpoint ----------------
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "TRL Research Backend is running 🚀",
		})
	})

	// ---------------- Admin ----------------
	admin := r.Group("/admins")
	{
		admin.POST("/", handlers.CreateAdmin)
		admin.GET("/:id", handlers.GetAdmin)
		admin.DELETE("/:id", handlers.DeleteAdmin)
	}

	// ---------------- Researcher ----------------
	researcher := r.Group("/researchers")
	{
		researcher.POST("/", handlers.CreateOrUpdateResearcher)
		researcher.GET("/:id", handlers.GetResearcher)
		researcher.DELETE("/:id", handlers.DeleteResearcher)
	}

	// ---------------- Case ----------------
	caseGroup := r.Group("/cases")
	{
		caseGroup.POST("/", handlers.CreateOrUpdateCase)
		caseGroup.GET("/:id", handlers.GetCase)
		caseGroup.DELETE("/:id", handlers.DeleteCase)
	}

	// ---------------- Coordinator ----------------
	coordinator := r.Group("/coordinators")
	{
		coordinator.POST("/", handlers.CreateOrUpdateCoordinator)
		coordinator.GET("/:email", handlers.GetCoordinator)
		coordinator.DELETE("/:email", handlers.DeleteCoordinator)
	}

	// ---------------- Intellectual Property ----------------
	ip := r.Group("/ips")
	{
		ip.POST("/", handlers.CreateOrUpdateIP)
		ip.GET("/:id", handlers.GetIP)
		ip.DELETE("/:id", handlers.DeleteIP)
	}

	// ---------------- Supporter ----------------
	supporter := r.Group("/supporters")
	{
		supporter.POST("/", handlers.CreateOrUpdateSupporter)
		supporter.GET("/:case_id", handlers.GetSupporter)
		supporter.DELETE("/:case_id", handlers.DeleteSupporter)
	}

	// ---------------- Assessment TRL ----------------
	trl := r.Group("/assessments")
	{
		trl.POST("/", handlers.CreateOrUpdateTrl)
		trl.GET("/:case_id", handlers.GetTrl)
		trl.DELETE("/:case_id", handlers.DeleteTrl)
	}

	// ---------------- Assessment TRL Part 1 ----------------
	trl1 := r.Group("/trl_part1")
	{
		trl1.POST("/", handlers.CreateOrUpdateTrlPart1)
		trl1.GET("/:id", handlers.GetTrlPart1)
		trl1.DELETE("/:id", handlers.DeleteTrlPart1)
	}

	// ---------------- Assessment TRL Part 2 ----------------
	trl2 := r.Group("/trl_part2")
	{
		trl2.POST("/", handlers.CreateOrUpdateTrlPart2)
		trl2.GET("/:id", handlers.GetTrlPart2)
		trl2.DELETE("/:id", handlers.DeleteTrlPart2)
	}

	// ---------------- Appointment ----------------
	appointment := r.Group("/appointments")
	{
		appointment.POST("/", handlers.CreateOrUpdateAppointment)
		appointment.GET("/:id", handlers.GetAppointment)
		appointment.DELETE("/:id", handlers.DeleteAppointment)
	}

	// ---------------- TRL Question (รวม part1+part2) ----------------
	/*
	trlq := r.Group("/trl_questions")
	{
		trlq.POST("/", handlers.CreateOrUpdateTrlQuestion)
		trlq.GET("/:assessment_id/:question_code", handlers.GetTrlQuestion)
		trlq.GET("/:assessment_id", handlers.GetTrlQuestionsByAssessment)
		trlq.DELETE("/:assessment_id/:question_code", handlers.DeleteTrlQuestion)
	}
	*/

	return r
}
