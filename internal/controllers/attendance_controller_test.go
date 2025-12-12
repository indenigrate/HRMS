package controllers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"hrms_backend/internal/controllers"
	"hrms_backend/internal/viewmodels"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Service ---
type MockAttendanceService struct {
	mock.Mock
}

func (m *MockAttendanceService) MarkAttendance(req viewmodels.CreateAttendanceRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAttendanceService) GetAttendanceByStudentID(studentID uint) ([]viewmodels.AttendanceResponse, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]viewmodels.AttendanceResponse), args.Error(1)
}

func (m *MockAttendanceService) GetWeeklyAttendance() ([]viewmodels.AttendanceResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]viewmodels.AttendanceResponse), args.Error(1)
}

// --- Tests ---

func TestMarkAttendanceController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAttendanceService)
	ctl := controllers.NewAttendanceController(mockService)
	r := gin.Default()
	r.POST("/attendance", ctl.MarkAttendance)

	// Case 1: Success
	mockService.On("MarkAttendance", mock.Anything).Return(nil).Once()

	// Need a valid ISO8601 date string
	reqBody := []byte(`{"student_id": 1, "date": "2025-12-12T09:00:00Z", "status": "present"}`)
	req, _ := http.NewRequest("POST", "/attendance", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Service Error (e.g., student not found)
	mockService.On("MarkAttendance", mock.Anything).Return(errors.New("student not found")).Once()
	req, _ = http.NewRequest("POST", "/attendance", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAttendanceByStudentController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAttendanceService)
	ctl := controllers.NewAttendanceController(mockService)
	r := gin.Default()
	r.GET("/attendance/student/:studentID", ctl.GetAttendanceByStudentID)

	// Case 1: Success
	expected := []viewmodels.AttendanceResponse{
		{ID: 1, StudentID: 1, Status: "present"},
	}
	mockService.On("GetAttendanceByStudentID", uint(1)).Return(expected, nil).Once()

	req, _ := http.NewRequest("GET", "/attendance/student/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "present")

	// Case 2: Invalid ID
	req, _ = http.NewRequest("GET", "/attendance/student/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
