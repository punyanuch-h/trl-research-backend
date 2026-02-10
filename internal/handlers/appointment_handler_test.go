package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trl-research-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppointmentRepository is a manual mock for AppointmentRepository
type MockAppointmentRepository struct {
	mock.Mock
}

func (m *MockAppointmentRepository) GetAppointmentAll() ([]models.Appointments, error) {
	args := m.Called()
	return args.Get(0).([]models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentByID(id string) (*models.Appointments, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentWithDetails(id string) (*models.Appointments, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentByCaseID(id string) ([]models.Appointments, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) CreateAppointment(ap *models.Appointments) error {
	args := m.Called(ap)
	return args.Error(0)
}

func (m *MockAppointmentRepository) UpdateAppointmentByID(id string, data map[string]interface{}) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockAppointmentRepository) GetUpcomingAppointments(start, end time.Time) ([]models.Appointments, error) {
	args := m.Called(start, end)
	return args.Get(0).([]models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) GetAppointmentByResearcherID(id string) ([]models.Appointments, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) UpdateNotifyStatus(id string, isNotify bool) error {
	args := m.Called(id, isNotify)
	return args.Error(0)
}

func (m *MockAppointmentRepository) GetNotificationsByRole(role string, userID string) ([]models.Appointments, error) {
	args := m.Called(role, userID)
	return args.Get(0).([]models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) GetUnreadNotificationCountByRole(role string, userID string) (int64, error) {
	args := m.Called(role, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAppointmentRepository) MarkNotificationAsRead(id string, role string, userID string) (*models.Appointments, error) {
	args := m.Called(id, role, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointments), args.Error(1)
}

func (m *MockAppointmentRepository) MarkAllNotificationsAsRead(role string, userID string) error {
	args := m.Called(role, userID)
	return args.Error(0)
}

func TestGetAppointmentByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockAppointmentRepository)
		handler := &AppointmentHandler{Repo: mockRepo}

		expectedAp := &models.Appointments{ID: "AP-00001", Detail: "Test Appointment"}
		mockRepo.On("GetAppointmentByID", "AP-00001").Return(expectedAp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "AP-00001"}}

		handler.GetAppointmentByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var got models.Appointments
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expectedAp.ID, got.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := new(MockAppointmentRepository)
		handler := &AppointmentHandler{Repo: mockRepo}

		mockRepo.On("GetAppointmentByID", "non-existent").Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "non-existent"}}

		handler.GetAppointmentByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockRepo.AssertExpectations(t)
	})
}
