package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockRepo       func(m *MockRefreshTokenRepository)
		expectedStatus int
	}{
		{
			name: "Happy path: Logout",
			requestBody: LogoutRequest{
				RefreshToken: "token-to-logout",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				hashed := utils.HashToken("token-to-logout")
				m.On("GetByHash", hashed).Return(MockRefreshTokenModel("rt-1"), nil)
				m.On("Revoke", "rt-1").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error: Token not found (still 200 for logout usually)",
			requestBody: LogoutRequest{
				RefreshToken: "missing",
			},
			mockRepo: func(m *MockRefreshTokenRepository) {
				m.On("GetByHash", mock.Anything).Return(nil, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRefreshTokenRepository)
			tt.mockRepo(repo)

			handler := &LogoutHandler{
				RefreshTokenRepo: repo,
			}

			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/auth/logout", bytes.NewBuffer(body))

			handler.Logout(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func MockRefreshTokenModel(id string) *models.RefreshToken {
	return &models.RefreshToken{
		ID: id,
	}
}
