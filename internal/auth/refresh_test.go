package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	mockKP := utils.NewManualKeyProvider("v1", privateKey, publicKey)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockRepo       func(m *MockRefreshTokenRepository)
		expectedStatus int
	}{
		{
			name: "Happy path: Refresh token",
			requestBody: RefreshRequest{
				RefreshToken: "valid-token",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				hashed := utils.HashToken("valid-token")
				m.On("GetByHash", hashed).Return(&models.RefreshToken{
					ID:        "rt-123",
					UserID:    "user-1",
					UserType:  "admin",
					ExpiryAt:  time.Now().Add(time.Hour),
					TokenHash: hashed,
				}, nil)
				m.On("Update", mock.Anything).Return(nil)
				m.On("Create", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error: Token expired",
			requestBody: RefreshRequest{
				RefreshToken: "expired-token",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				m.On("GetByHash", mock.Anything).Return(&models.RefreshToken{
					ExpiryAt: time.Now().Add(-time.Hour),
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Error: Token revoked",
			requestBody: RefreshRequest{
				RefreshToken: "revoked-token",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				now := time.Now()
				m.On("GetByHash", mock.Anything).Return(&models.RefreshToken{
					RevokedAt: &now,
					ExpiryAt:  time.Now().Add(time.Hour),
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Security: Reuse detection",
			requestBody: RefreshRequest{
				RefreshToken: "reused-token",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				now := time.Now()
				m.On("GetByHash", mock.Anything).Return(&models.RefreshToken{
					UserID:          "victim-user",
					ReplacedByToken: "attacker-token",
					ExpiryAt:        time.Now().Add(time.Hour),
					RevokedAt:       &now,
				}, nil)
				m.On("RevokeAllForUser", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRefreshTokenRepository)
			adminRepo := new(MockAdminRepository)
			researcherRepo := new(MockResearcherRepository)

			tt.mockRepo(repo)

			// Default mocks for user lookup if needed
			adminRepo.On("GetAdminByID", mock.Anything).Return(&models.Admins{Email: "admin@test.com"}, nil)
			researcherRepo.On("GetResearcherByID", mock.Anything).Return(&models.Researchers{Email: "res@test.com"}, nil)

			handler := &RefreshHandler{
				RefreshTokenRepo: repo,
				AdminRepo:        adminRepo,
				ResearcherRepo:   researcherRepo,
				KeyProvider:      mockKP,
				Cfg: config.Config{
					JWTExpiry:          "15",
					RefreshTokenExpiry: "10h",
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))

			handler.Refresh(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
