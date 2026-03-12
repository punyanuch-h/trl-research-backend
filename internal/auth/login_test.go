package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"trl-research-backend/internal/config"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/repository"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup Test Keys
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	mockKP := utils.NewManualKeyProvider("v1", privateKey, publicKey)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockAdmin      func(m *MockAdminRepository)
		mockResearcher func(m *MockResearcherRepository)
		expectedStatus int
		expectedRole   string
	}{
		{
			name: "Happy path: Admin login",
			requestBody: LoginRequest{
				Email:    "admin@test.com",
				Password: "password123",
			},
			mockAdmin: func(m *MockAdminRepository) {
				m.On("Login", "admin@test.com", "password123").Return(&models.Admins{ID: "AD-001", Email: "admin@test.com"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedRole:   "admin",
		},
		{
			name: "Happy path: Researcher login",
			requestBody: LoginRequest{
				Email:    "researcher@test.com",
				Password: "password123",
			},
			mockAdmin: func(m *MockAdminRepository) {
				m.On("Login", "researcher@test.com", "password123").Return(nil, errors.New("not found"))
			},
			mockResearcher: func(m *MockResearcherRepository) {
				m.On("Login", "researcher@test.com", "password123").Return(&models.Researchers{ID: "RS-001", Email: "researcher@test.com"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedRole:   "researcher",
		},
		{
			name: "Error: Invalid credentials",
			requestBody: LoginRequest{
				Email:    "wrong@test.com",
				Password: "wrong",
			},
			mockAdmin: func(m *MockAdminRepository) {
				m.On("Login", "wrong@test.com", "wrong").Return(nil, errors.New("unauthorized"))
			},
			mockResearcher: func(m *MockResearcherRepository) {
				m.On("Login", "wrong@test.com", "wrong").Return(nil, errors.New("unauthorized"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Error: Temporary password expired",
			requestBody: LoginRequest{
				Email:    "temp@test.com",
				Password: "password123",
			},
			mockAdmin: func(m *MockAdminRepository) {
				m.On("Login", "temp@test.com", "password123").Return(nil, repository.ErrTempPasswordExpired)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Error: Empty request",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminRepo := new(MockAdminRepository)
			researcherRepo := new(MockResearcherRepository)
			if tt.mockAdmin != nil {
				tt.mockAdmin(adminRepo)
			}
			if tt.mockResearcher != nil {
				tt.mockResearcher(researcherRepo)
			}
			refreshTokenRepo := new(MockRefreshTokenRepository)
			refreshTokenRepo.On("Create", mock.Anything).Return(nil)

			handler := &LoginHandler{
				AdminRepo:        adminRepo,
				ResearcherRepo:   researcherRepo,
				RefreshTokenRepo: refreshTokenRepo,
				KeyProvider:      mockKP,
				Cfg: config.Config{
					JWTExpiry:          "480",
					JWTExpiryTemp:      "10",
					RefreshTokenExpiry: "10h",
				},
			}

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))

			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, tt.expectedRole, resp["role"])
				assert.NotEmpty(t, resp["token"])
				assert.NotEmpty(t, resp["refresh_token"])
				assert.Contains(t, resp, "is_temp")
				assert.Contains(t, resp, "expires_in")
				assert.Equal(t, "minutes", resp["unit"])

				// Check types
				_, ok := resp["is_temp"].(bool)
				assert.True(t, ok, "is_temp should be a boolean")
				_, ok = resp["expires_in"].(float64) // JSON numbers are float64 in Go map[string]interface{}
				assert.True(t, ok, "expires_in should be a number")
			} else if tt.name == "Error: Temporary password expired" {
				var resp map[string]interface{}
				_ = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Equal(t, "temporary password expired", resp["error"])
			}
		})
	}
}
