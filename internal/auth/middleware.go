package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(header, "Bearer ")

        // Validate required environment variables
        jwtIssuer := os.Getenv("JWT_ISSUER")
        jwtAudience := os.Getenv("JWT_AUDIENCE")
        if jwtIssuer == "" {
            log.Println("❌ JWT_ISSUER environment variable is not set")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server configuration error: JWT_ISSUER not set"})
            c.Abort()
            return
        }
        if jwtAudience == "" {
            log.Println("❌ JWT_AUDIENCE environment variable is not set")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "server configuration error: JWT_AUDIENCE not set"})
            c.Abort()
            return
        }

        kp, err := utils.NewEnvKeyProvider()
        if err != nil {
            log.Printf("❌ Key provider error: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("key provider error: %v", err)})
            c.Abort()
            return
        }

        claims, err := utils.ValidateJWT(
            tokenString,
            jwtIssuer,
            jwtAudience,
            *kp,
        )
        if err != nil {
            log.Printf("❌ JWT validation failed: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("invalid token: %v", err)})
            c.Abort()
            return
        }

        log.Println("✅ AuthMiddleware: JWT claims loaded:")
		log.Printf("→ UserID: %s", claims.UserID)
		log.Printf("→ Email: %s", claims.UserEmail)
		log.Printf("→ Role: %s", claims.Role)

        // Inject claims into context
        c.Set("authCtx", claims)
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.UserEmail)
		c.Set("role", claims.Role)

        c.Next()
    }
}

func GetMiddleware(c *gin.Context) (*utils.Claims, error) {
	val, ok := c.Get("authCtx")
	if !ok {
		return nil, errors.New("auth context missing")
	}
	claims, ok := val.(*utils.Claims)
	if !ok {
		return nil, errors.New("invalid auth context type")
	}
	return claims, nil
}
