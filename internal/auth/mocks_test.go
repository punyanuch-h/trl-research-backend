package auth

import (
	"time"
	"trl-research-backend/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockRefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByHash(tokenHash string) (*models.RefreshToken, error) {
	args := m.Called(tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Update(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// MockAdminRepository
type MockAdminRepository struct {
	mock.Mock
}

func (m *MockAdminRepository) GetAdminAll() ([]models.Admins, error)                { return nil, nil }
func (m *MockAdminRepository) GetAdminByID(id string) (*models.Admins, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Admins), args.Error(1)
}
func (m *MockAdminRepository) GetAdminByEmail(email string) (*models.Admins, error) { return nil, nil }
func (m *MockAdminRepository) CreateAdmin(admin *models.Admins) error               { return nil }
func (m *MockAdminRepository) UpdatePasswordByEmail(email string, password string, isTemp bool, expiresAt *time.Time) error {
	return nil
}
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
func (m *MockResearcherRepository) UpdatePasswordByEmail(email string, password string, isTemp bool, expiresAt *time.Time) error {
	return nil
}
func (m *MockResearcherRepository) GetResearcherAll() ([]models.Researchers, error) { return nil, nil }
func (m *MockResearcherRepository) GetResearcherByID(id string) (*models.Researchers, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Researchers), args.Error(1)
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
