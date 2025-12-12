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
type MockStudentService struct {
	mock.Mock
}

func (m *MockStudentService) CreateStudent(req viewmodels.CreateStudentRequest) (*viewmodels.StudentResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*viewmodels.StudentResponse), args.Error(1)
}

func (m *MockStudentService) GetAllStudents(page, limit int) ([]viewmodels.StudentResponse, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]viewmodels.StudentResponse), args.Error(1)
}

func (m *MockStudentService) GetStudentByID(id uint) (*viewmodels.StudentResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*viewmodels.StudentResponse), args.Error(1)
}

func (m *MockStudentService) UpdateStudent(id uint, req viewmodels.UpdateStudentRequest) (*viewmodels.StudentResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*viewmodels.StudentResponse), args.Error(1)
}

func (m *MockStudentService) DeleteStudent(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Helper to setup router ---
func setupRouter(service *MockStudentService) (*controllers.StudentController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	ctl := controllers.NewStudentController(service)
	r := gin.Default()
	// Register routes manually for testing
	r.POST("/students", ctl.CreateStudent)
	r.GET("/students", ctl.GetAllStudents)
	r.GET("/students/:id", ctl.GetStudentByID)
	r.PUT("/students/:id", ctl.UpdateStudent)
	r.DELETE("/students/:id", ctl.DeleteStudent)
	return ctl, r
}

// --- Tests ---

func TestCreateStudentController(t *testing.T) {
	mockService := new(MockStudentService)
	_, r := setupRouter(mockService)

	// Case 1: Success
	expected := &viewmodels.StudentResponse{ID: 1, Name: "Alice"}
	mockService.On("CreateStudent", mock.Anything).Return(expected, nil).Once()

	reqBody := []byte(`{"name":"Alice","email":"a@a.com","department":"IT"}`)
	req, _ := http.NewRequest("POST", "/students", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Bad Request (Invalid JSON)
	reqBody = []byte(`{"name": ""}`) // Missing required email
	req, _ = http.NewRequest("POST", "/students", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllStudentsController(t *testing.T) {
	mockService := new(MockStudentService)
	_, r := setupRouter(mockService)

	// Case 1: Success
	expected := []viewmodels.StudentResponse{{ID: 1, Name: "Alice"}}
	mockService.On("GetAllStudents", 1, 10).Return(expected, nil).Once()

	req, _ := http.NewRequest("GET", "/students?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Alice")
}

func TestGetStudentByIDController(t *testing.T) {
	mockService := new(MockStudentService)
	_, r := setupRouter(mockService)

	// Case 1: Success
	expected := &viewmodels.StudentResponse{ID: 1, Name: "Alice"}
	mockService.On("GetStudentByID", uint(1)).Return(expected, nil).Once()

	req, _ := http.NewRequest("GET", "/students/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Not Found
	mockService.On("GetStudentByID", uint(99)).Return(nil, errors.New("not found")).Once()
	req, _ = http.NewRequest("GET", "/students/99", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Case 3: Invalid ID (String instead of int)
	req, _ = http.NewRequest("GET", "/students/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateStudentController(t *testing.T) {
	mockService := new(MockStudentService)
	_, r := setupRouter(mockService)

	// Case 1: Success
	expected := &viewmodels.StudentResponse{ID: 1, Name: "Updated"}
	mockService.On("UpdateStudent", uint(1), mock.Anything).Return(expected, nil).Once()

	reqBody := []byte(`{"name":"Updated"}`)
	req, _ := http.NewRequest("PUT", "/students/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Case 2: Invalid ID
	req, _ = http.NewRequest("PUT", "/students/abc", bytes.NewBuffer(reqBody))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteStudentController(t *testing.T) {
	mockService := new(MockStudentService)
	_, r := setupRouter(mockService)

	// Case 1: Success
	mockService.On("DeleteStudent", uint(1)).Return(nil).Once()
	req, _ := http.NewRequest("DELETE", "/students/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Case 2: Service Error
	mockService.On("DeleteStudent", uint(99)).Return(errors.New("failed")).Once()
	req, _ = http.NewRequest("DELETE", "/students/99", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
