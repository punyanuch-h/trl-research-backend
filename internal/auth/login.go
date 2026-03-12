package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// LoginHandler รวม repository ของทั้ง admin และ researcher
type LoginHandler struct {
	AdminRepo        repository.AdminRepository
	ResearcherRepo   repository.ResearcherRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	KeyProvider      utils.IKeyProvider
	Cfg              config.Config
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
	if errors.Is(errA, repository.ErrTempPasswordExpired) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "temporary password expired"})
		return
	}
	if errA == nil && admin != nil {
		userID = admin.ID
		userEmail = admin.Email
		userRole = "admin"
	}

	// 2️⃣ ถ้ายังไม่เจอใน admin ให้ลองเช็ก researcher
	var researcher *models.Researchers
	if userRole == "" {
		var errR error
		researcher, errR = h.ResearcherRepo.Login(req.Email, req.Password)
		fmt.Println("researcher", researcher)
		fmt.Println("errR", errR)
		if errors.Is(errR, repository.ErrTempPasswordExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "temporary password expired"})
			return
		}
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

	// 4️⃣ โหลด key provider
	kp := h.KeyProvider
	if kp == nil {
		var err error
		kp, err = utils.NewEnvKeyProvider()
		if err != nil {
			fmt.Println("❌ key provider init failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal key provider error"})
			return
		}
	}

	// 5️⃣ สร้าง JWT token
	// ตัดสินใจเลือก expiry ตามสถานะ PasswordIsTemp
	var isTemp bool
	if admin != nil {
		isTemp = admin.PasswordIsTemp
		fmt.Println("isTemp", isTemp)
	} else if researcher != nil {
		isTemp = researcher.PasswordIsTemp
		fmt.Println("isTemp", isTemp)
	}

	var ttl time.Duration
	if isTemp {
		ttl = h.Cfg.GetJWTExpiryTemp()
	} else {
		ttl = h.Cfg.GetJWTExpiry()
	}

	if ttl <= 0 {
		fmt.Println("❌ invalid ttl detected:", ttl)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token expiry configuration"})
		return
	}

	token, err := utils.GenerateJWT(
		userID,    // user id
		userEmail, // email
		userRole,  // role (admin/researcher)
		"", "",    // clientID, clientName (optional)
		os.Getenv("JWT_ISSUER"),
		os.Getenv("JWT_AUDIENCE"),
		"v1", // key id
		isTemp,
		ttl,
		kp,
	)
	if err != nil {
		fmt.Println("❌ failed to generate token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	// 5.1️⃣ Generate Refresh Token
	refreshTokenStr, err := utils.GenerateRandomToken()
	if err != nil {
		fmt.Println("❌ failed to generate refresh token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate refresh token"})
		return
	}

	refreshExpiry := h.Cfg.GetRefreshTokenExpiry()
	refreshTokenModel := &models.RefreshToken{
		UserID:    userID,
		TokenHash: utils.HashToken(refreshTokenStr),
		ExpiryAt:  time.Now().Add(refreshExpiry),
		UserType:  userRole,
	}

	if err := h.RefreshTokenRepo.Create(refreshTokenModel); err != nil {
		fmt.Println("❌ failed to save refresh token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal storage error"})
		return
	}

	// 6️⃣ ส่ง response กลับไป
	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshTokenStr,
		"expires_in":    int(ttl.Minutes()),
		"unit":          "minutes",
		"role":          userRole,
		"is_temp":       isTemp,
	})
}
