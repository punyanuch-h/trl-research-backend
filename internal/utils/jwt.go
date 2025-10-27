package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	Role             string `json:"role"`
	ClientID         string `json:"client_id"`
	ClientName       string `json:"client_name"`
	jwt.RegisteredClaims
}

// GenerateJWT
func GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid string, ttl time.Duration, kp KeyProvider) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID:    userID,
		UserEmail: userEmail,
		Role:             role,
		ClientID:         clientID,
		ClientName:       clientName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  []string{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-30 * time.Second)), // leeway
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl * time.Hour)),
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
func ValidateJWT(tokenString, expectedIssuer, expectedAudience string, kp KeyProvider) (*Claims, error) {
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
