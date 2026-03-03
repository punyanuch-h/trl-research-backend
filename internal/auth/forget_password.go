package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"
	"trl-research-backend/internal/utils/send_email"

	"github.com/gin-gonic/gin"
)

type ForgotHandler struct {
	AdminRepo      repository.AdminRepository
	ResearcherRepo repository.ResearcherRepository
	Cfg            config.Config
}

type ForgotReq struct {
	Email string `json:"email"`
}

func (h *ForgotHandler) ForgetPassword(c *gin.Context) {
	var req ForgotReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		if err != nil {
			log.Printf("ForgetPassword Bind Error: %v", err)
		} else {
			log.Println("ForgetPassword Error: Email is empty")
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	var userRole string

	// 1. Check if user is an Admin
	_, err := h.AdminRepo.GetAdminByEmail(req.Email)
	if err == nil {
		userRole = "admin"
	}

	// 2. If not admin, check if user is a Researcher
	if userRole == "" {
		_, err := h.ResearcherRepo.GetResearcherByEmail(req.Email)
		if err == nil {
			userRole = "researcher"
		}
	}

	// 3. If no user found in either, return success to prevent email enumeration
	if userRole == "" {
		c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
		return
	}

	// 4. Generate temp password
	tempPass, err := utils.GenerateTempPassword(24)
	if err != nil {
		log.Printf("ForgotPassword Error: failed to generate password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// 5. Update password in the correct repo
	expiresAt := time.Now().Add(h.Cfg.GetTempPasswordExpiry())
	if userRole == "admin" {
		if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, string(tempPass), true, &expiresAt); err != nil {
			log.Printf("ForgotPassword Error (Admin): %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	} else if userRole == "researcher" {
		if err := h.ResearcherRepo.UpdatePasswordByEmail(req.Email, string(tempPass), true, &expiresAt); err != nil {
			log.Printf("ForgotPassword Error (Researcher): %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
	}

	// 6. Send email
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		log.Printf("ForgotPassword Error: invalid EMAIL_PORT configuration: %v", portStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	sender := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	if host == "" || sender == "" || pass == "" {
		log.Printf("ForgotPassword Error: missing email SMTP configuration")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	smtpService := send_email.CreateSMTPEmailService(send_email.SMTPConfig{
		Host:     host,
		Port:     portStr,
		Username: sender,
		Password: pass,
		From:     sender,
	})

	body := send_email.TemplateForgetPassword(string(tempPass))
	if err := smtpService(req.Email, "Password Reset Request", body); err != nil {
		log.Printf("ForgotPassword Error: SMTP failed: %v", err)
		// We could return success here too to avoid leaking info, but for debugging let's keep it
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
}
