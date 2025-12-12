package services_test

import (
	"errors"
	"testing"
	"time"

	"hrms_backend/internal/models"
	"hrms_backend/internal/services"
	"hrms_backend/internal/viewmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mock Attendance Repo ---
type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) Create(attendance *models.Attendance) error {
	args := m.Called(attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetAttendanceByStudentID(studentID uint) ([]models.Attendance, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Attendance), args.Error(1)
}

func (m *MockAttendanceRepo) GetAttendanceSince(date time.Time) ([]models.Attendance, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Attendance), args.Error(1)
}

// --- Tests ---

func TestMarkAttendance(t *testing.T) {
	mockAttRepo := new(MockAttendanceRepo)
	mockStudentRepo := new(MockStudentRepo) // Reusing the mock from student_service_test.go
	service := services.NewAttendanceService(mockAttRepo, mockStudentRepo)

	req := viewmodels.CreateAttendanceRequest{StudentID: 1, Date: time.Now(), Status: "present"}

	// Case 1: Success
	// Expect check for student existence first
	mockStudentRepo.On("GetByID", uint(1)).Return(&models.Student{Model: gorm.Model{ID: 1}}, nil).Once()
	// Then expect create attendance
	mockAttRepo.On("Create", mock.AnythingOfType("*models.Attendance")).Return(nil).Once()

	err := service.MarkAttendance(req)
	assert.NoError(t, err)

	// Case 2: Student Not Found
	mockStudentRepo.On("GetByID", uint(1)).Return(nil, errors.New("not found")).Once()
	err = service.MarkAttendance(req)
	assert.Error(t, err)
	assert.Equal(t, "student not found: cannot mark attendance", err.Error())
}

func TestGetAttendanceByStudentID(t *testing.T) {
	mockAttRepo := new(MockAttendanceRepo)
	mockStudentRepo := new(MockStudentRepo)
	service := services.NewAttendanceService(mockAttRepo, mockStudentRepo)

	// Case 1: Success
	mockStudentRepo.On("GetByID", uint(1)).Return(&models.Student{}, nil).Once()
	mockData := []models.Attendance{
		{Model: gorm.Model{ID: 1}, StudentID: 1, Status: "present"},
	}
	mockAttRepo.On("GetAttendanceByStudentID", uint(1)).Return(mockData, nil).Once()

	resp, err := service.GetAttendanceByStudentID(1)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	// Case 2: Student Not Found
	mockStudentRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found")).Once()
	resp, err = service.GetAttendanceByStudentID(99)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetWeeklyAttendance(t *testing.T) {
	mockAttRepo := new(MockAttendanceRepo)
	mockStudentRepo := new(MockStudentRepo)
	service := services.NewAttendanceService(mockAttRepo, mockStudentRepo)

	// Case 1: Success
	mockData := []models.Attendance{
		{Model: gorm.Model{ID: 1}, StudentID: 1, Status: "present"},
	}
	// We use mock.Anything for the date argument since exact time matching is flaky
	mockAttRepo.On("GetAttendanceSince", mock.Anything).Return(mockData, nil).Once()

	resp, err := service.GetWeeklyAttendance()
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
}
