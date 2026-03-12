package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string `json:"user_id"`
	UserEmail  string `json:"user_email"`
	Role       string `json:"role"`
	ClientID   string `json:"client_id"`
	ClientName string `json:"client_name"`
	IsTemp     bool   `json:"is_temp"`
	jwt.RegisteredClaims
}

// GenerateJWT
func GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid string, isTemp bool, ttl time.Duration, kp IKeyProvider) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:     userID,
		UserEmail:  userEmail,
		Role:       role,
		ClientID:   clientID,
		ClientName: clientName,
		IsTemp:     isTemp,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)), // leeway
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	privateKey, err := kp.GetPrivateKey(kid)
	if err != nil {
		return "", err
	}
	return token.SignedString(privateKey)
}

// ValidateJWT
func ValidateJWT(tokenString, expectedIssuer, expectedAudience string, kp IKeyProvider) (*Claims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithAudience(expectedAudience),
		jwt.WithLeeway(30*time.Second),
	)

	claims := &Claims{}
	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid")
		}
		return kp.GetPublicKey(kid)
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// GenerateRandomToken generates a secure random string for use as a Refresh Token
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashToken hashes a token string using SHA-256
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
