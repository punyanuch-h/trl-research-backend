package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginHandler รวม repository ของทั้ง admin และ researcher
type LoginHandler struct {
	AdminRepo      repository.AdminRepository
	ResearcherRepo repository.ResearcherRepository
}

// LoginRequest รับข้อมูลจาก frontend
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login ฟังก์ชันหลักในการตรวจสอบผู้ใช้และสร้าง JWT token
func (h *LoginHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}


	var userID, userEmail, userRole string

	// 1️⃣ ตรวจสอบในตาราง Admin ก่อน
	admin, errA := h.AdminRepo.Login(req.Email, req.Password)
	fmt.Println("admin", admin)
	fmt.Println("errA", errA)
	if errA == nil && admin != nil {
		userID = admin.ID
		userEmail = admin.Email
		userRole = "admin"
	}

	// 2️⃣ ถ้ายังไม่เจอใน admin ให้ลองเช็ก researcher
	if userRole == "" {
		researcher, errR := h.ResearcherRepo.Login(req.Email, req.Password)
		fmt.Println("researcher", researcher)
		fmt.Println("errR", errR)
		if errR == nil && researcher != nil {
			userID = researcher.ID
			userEmail = researcher.Email
			userRole = "researcher"
		}
	}

	// 3️⃣ ถ้าไม่เจอทั้งสองกลุ่ม
	if userRole == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	fmt.Printf("✅ Verified user: %s (role: %s)\n", userEmail, userRole)

	// 4️⃣ โหลด key provider จาก environment
	kp, err := utils.NewEnvKeyProvider()
	if err != nil {
		fmt.Println("❌ key provider init failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal key provider error"})
		return
	}

	// 5️⃣ สร้าง JWT token
	expH, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if expH <= 0 {
		expH = 24
	}

	token, err := utils.GenerateJWT(
		userID,    // user id
		userEmail, // email
		userRole,  // role (admin/researcher)
		"", "",    // clientID, clientName (optional)
		os.Getenv("JWT_ISSUER"),
		os.Getenv("JWT_AUDIENCE"),
		"v1", // key id
		time.Duration(expH),
		*kp,
	)
	if err != nil {
		fmt.Println("❌ failed to generate token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	// 6️⃣ ส่ง response กลับไป
	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": expH,
		"role":       userRole,
	})
}
