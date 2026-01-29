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

func (h *ForgotHandler) ForgotPassword(c *gin.Context) {
	var req ForgotReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		if err != nil {
			log.Printf("ForgotPassword Bind Error: %v", err)
		} else {
			log.Println("ForgotPassword Error: Email is empty")
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": fmt.Sprintf("bind_err: %v, email: %s", err, req.Email),
		})
		return
	}

	// prevent enumeration
	if _, err := h.AdminRepo.GetAdminByEmail(req.Email); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
		return
	}
	fmt.Println("check admins")

	// temp password
	tempPass := utils.GenerateTempPassword(24)
	fmt.Println("tempPass", tempPass)
	if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, string(tempPass)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	fmt.Println("update password success")

	// send email
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	user := os.Getenv("EMAIL_SENDER")
	pass := os.Getenv("EMAIL_PASSWORD")

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", user, pass, host)
	msg := []byte("Subject: Temporary Password\r\n\r\nYour temporary password is: " + tempPass)
	if err := smtp.SendMail(addr, auth, user, []string{req.Email}, msg); err != nil {
		// prevent leak
		fmt.Println("send email error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	fmt.Println("send email success")

	c.JSON(http.StatusOK, gin.H{"message": "Temporary password has been sent"})
}
