package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKey := &privateKey.PublicKey

	kid := "test-key"
	kp := NewManualKeyProvider(kid, privateKey, publicKey)

	userID := "user-123"
	userEmail := "test@example.com"
	role := "admin"
	clientID := "client-abc"
	clientName := "Test Client"
	issuer := "test-issuer"
	audience := "test-audience"
	ttl := time.Hour // 1 hour

	t.Run("Generate and Validate JWT", func(t *testing.T) {
		token, err := GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid, false, ttl, kp)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := ValidateJWT(token, issuer, audience, kp)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, userEmail, claims.UserEmail)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, clientID, claims.ClientID)
		assert.Equal(t, clientName, claims.ClientName)
		assert.Equal(t, issuer, claims.Issuer)
		assert.Contains(t, claims.Audience, audience)
	})

	t.Run("Validate JWT with wrong issuer", func(t *testing.T) {
		token, err := GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid, false, ttl, kp)
		assert.NoError(t, err)

		_, err = ValidateJWT(token, "wrong-issuer", audience, kp)
		assert.Error(t, err)
	})

	t.Run("Validate JWT with wrong audience", func(t *testing.T) {
		token, err := GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid, false, ttl, kp)
		assert.NoError(t, err)

		_, err = ValidateJWT(token, issuer, "wrong-audience", kp)
		assert.Error(t, err)
	})

	t.Run("Validate JWT with wrong kid", func(t *testing.T) {
		token, err := GenerateJWT(userID, userEmail, role, clientID, clientName, issuer, audience, kid, false, ttl, kp)
		assert.NoError(t, err)

		kpWrong := &KeyProvider{
			privateKeys: map[string]*rsa.PrivateKey{},
			publicKeys:  map[string]*rsa.PublicKey{},
		}
		_, err = ValidateJWT(token, issuer, audience, kpWrong)
		assert.Error(t, err)
	})
}

func TestRandomTokenAndHash(t *testing.T) {
	t.Run("Generate Random Token", func(t *testing.T) {
		token1, err := GenerateRandomToken()
		assert.NoError(t, err)
		assert.Len(t, token1, 64) // 32 bytes hex encoded

		token2, err := GenerateRandomToken()
		assert.NoError(t, err)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("Hash Token", func(t *testing.T) {
		token := "my-secret-token"
		hash1 := HashToken(token)
		assert.NotEmpty(t, hash1)

		hash2 := HashToken(token)
		assert.Equal(t, hash1, hash2)

		differentToken := "another-token"
		hash3 := HashToken(differentToken)
		assert.NotEqual(t, hash1, hash3)
	})
}
