package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"trl-research-backend/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ForgotHandler struct {
	AdminRepo repository.AdminRepo
}

type ForgotReq struct {
	Email string `json:"email"`
}

func (h *ForgotHandler) ForgotPassword(c *gin.Context) {
	var req ForgotReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ตรวจว่ามี user ไหม
	if _, err := h.AdminRepo.GetAdminByEmail(req.Email); err != nil {
		// ป้องกัน enumeration: ตอบ 200 แต่ไม่บอกว่าไม่มี
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a temporary password has been sent"})
		return
	}

	// รหัสชั่วคราว
	rand.Seed(time.Now().UnixNano())
	tempPass := fmt.Sprintf("%06d", rand.Intn(1000000))

	hash, err := bcrypt.GenerateFromPassword([]byte(tempPass), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if err := h.AdminRepo.UpdatePasswordByEmail(req.Email, string(hash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// ส่งอีเมล (ตัวอย่าง SMTP พื้นฐาน; คุณใช้ gomail หรือบริการอื่นแทนได้)
	host := os.Getenv("EMAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	user := os.Getenv("EMAIL_SENDER")
	pass := os.Getenv("EMAIL_PASSWORD")

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", user, pass, host)
	msg := []byte("Subject: Temporary Password\r\n\r\nYour temporary password is: " + tempPass)
	if err := smtp.SendMail(addr, auth, user, []string{req.Email}, msg); err != nil {
		// ถึงส่งเมลพลาด ก็ไม่ leak ให้ attacker รู้
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a temporary password has been sent"})
}
