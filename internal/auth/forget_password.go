package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type ForgotHandler struct {
	AdminRepo repository.AdminRepository
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

	// prevent enumeration
	if _, err := h.AdminRepo.GetAdminByEmail(req.Email); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
		return
	}

	if _,err := h.ResearcherRepo.GetResearcherByEmail(req.Email); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
		return
	}

	// temp password
	tempPass, err := utils.GenerateTempPassword(24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, string(tempPass)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// send email
	host := os.Getenv("EMAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil || port <= 0 {
		log.Printf("ForgotPassword Error: invalid EMAIL_PORT configuration")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	user := os.Getenv("EMAIL_SENDER")
	pass := os.Getenv("EMAIL_PASSWORD")

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", user, pass, host)
	msg := []byte("Subject: Temporary Password\r\n\r\nYour temporary password is: " + tempPass)
	fmt.Println(msg)
	if err := smtp.SendMail(addr, auth, user, []string{req.Email}, msg); err != nil {
		// prevent leak
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
}
