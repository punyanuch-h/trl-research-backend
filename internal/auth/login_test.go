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
	"trl-research-backend/internal/models"
	"trl-research-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) GetAdminAll() ([]models.Admins, error)                     { return nil, nil }
func (m *MockAdminRepository) GetAdminByID(id string) (*models.Admins, error)            { return nil, nil }
func (m *MockAdminRepository) GetAdminByEmail(email string) (*models.Admins, error)      { return nil, nil }
func (m *MockAdminRepository) CreateAdmin(admin *models.Admins) error                    { return nil }
func (m *MockAdminRepository) UpdatePasswordByEmail(email string, password string) error { return nil }
func (m *MockAdminRepository) UpdateAdminByID(adminID string, data *models.Admins) error { return nil }
func (m *MockAdminRepository) DeleteAdmin(email string) error                            { return nil }

func (m *MockAdminRepository) Login(email, password string) (*models.Admins, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admins), args.Error(1)
}

// MockResearcherRepository
type MockResearcherRepository struct {
	mock.Mock
}

func (m *MockResearcherRepository) GetResearcherByEmail(email string) (*models.Researchers, error) {
	return nil, nil
}
func (m *MockResearcherRepository) UpdatePasswordByEmail(email string, password string) error {
	return nil
}
func (m *MockResearcherRepository) GetResearcherAll() ([]models.Researchers, error) { return nil, nil }
func (m *MockResearcherRepository) GetResearcherByID(id string) (*models.Researchers, error) {
	return nil, nil
}
func (m *MockResearcherRepository) GetResearcherByCaseID(id string) (*models.Researchers, error) {
	return nil, nil
}
func (m *MockResearcherRepository) CreateResearcher(r *models.Researchers) error { return nil }
func (m *MockResearcherRepository) UpdateResearcherByID(id string, data *models.Researchers) error {
	return nil
}

func (m *MockResearcherRepository) Login(email, password string) (*models.Researchers, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Researchers), args.Error(1)
}

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

			handler := &LoginHandler{
				AdminRepo:      adminRepo,
				ResearcherRepo: researcherRepo,
				KeyProvider:    mockKP,
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
			}
		})
	}
}
