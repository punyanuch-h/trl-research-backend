package auth

import (
	"errors"
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

        kp, err := utils.NewEnvKeyProvider()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "key provider error"})
            c.Abort()
            return
        }

        claims, err := utils.ValidateJWT(
            tokenString,
            os.Getenv("JWT_ISSUER"),
            os.Getenv("JWT_AUDIENCE"),
            *kp,
        )
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        // Inject claims into context
        c.Set("authCtx", claims)
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
