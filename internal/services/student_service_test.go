package services_test

import (
	"errors"
	"testing"

	"hrms_backend/internal/models"
	"hrms_backend/internal/services"
	"hrms_backend/internal/viewmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mock Repository ---
type MockStudentRepo struct {
	mock.Mock
}

func (m *MockStudentRepo) Create(student *models.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepo) GetAll(limit, offset int) ([]models.Student, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepo) GetByID(id uint) (*models.Student, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepo) Update(id uint, student *models.Student) error {
	args := m.Called(id, student)
	return args.Error(0)
}

func (m *MockStudentRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Tests ---

func TestCreateStudent(t *testing.T) {
	mockRepo := new(MockStudentRepo)
	service := services.NewStudentService(mockRepo)

	req := viewmodels.CreateStudentRequest{Name: "Alice", Email: "alice@test.com", Department: "IT"}

	// Case 1: Success
	mockRepo.On("Create", mock.AnythingOfType("*models.Student")).Return(nil).Once()
	resp, err := service.CreateStudent(req)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", resp.Name)

	// Case 2: DB Error
	mockRepo.On("Create", mock.Anything).Return(errors.New("db error")).Once()
	resp, err = service.CreateStudent(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetAllStudents(t *testing.T) {
	mockRepo := new(MockStudentRepo)
	service := services.NewStudentService(mockRepo)

	mockData := []models.Student{
		{Model: gorm.Model{ID: 1}, Name: "A", Email: "a@a.com"},
		{Model: gorm.Model{ID: 2}, Name: "B", Email: "b@b.com"},
	}

	// Case 1: Success with Pagination (Page 1, Limit 10 -> Offset 0)
	mockRepo.On("GetAll", 10, 0).Return(mockData, nil).Once()
	resp, err := service.GetAllStudents(1, 10)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)

	// Case 2: DB Error
	mockRepo.On("GetAll", 10, 0).Return(nil, errors.New("db error")).Once()
	resp, err = service.GetAllStudents(1, 10)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetStudentByID(t *testing.T) {
	mockRepo := new(MockStudentRepo)
	service := services.NewStudentService(mockRepo)

	student := &models.Student{Model: gorm.Model{ID: 1}, Name: "Alice"}

	// Case 1: Found
	mockRepo.On("GetByID", uint(1)).Return(student, nil).Once()
	resp, err := service.GetStudentByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", resp.Name)

	// Case 2: Not Found
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found")).Once()
	resp, err = service.GetStudentByID(99)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateStudent(t *testing.T) {
	mockRepo := new(MockStudentRepo)
	service := services.NewStudentService(mockRepo)

	existing := &models.Student{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	req := viewmodels.UpdateStudentRequest{Name: "New Name"}

	// Case 1: Success
	// Sequence: GetByID (Check existence) -> Update (Save) -> GetByID (Fetch updated)
	mockRepo.On("GetByID", uint(1)).Return(existing, nil).Once()
	mockRepo.On("Update", uint(1), mock.MatchedBy(func(s *models.Student) bool {
		return s.Name == "New Name"
	})).Return(nil).Once()
	mockRepo.On("GetByID", uint(1)).Return(&models.Student{Model: gorm.Model{ID: 1}, Name: "New Name"}, nil).Once()

	resp, err := service.UpdateStudent(1, req)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", resp.Name)

	// Case 2: Student Not Found
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found")).Once()
	resp, err = service.UpdateStudent(99, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteStudent(t *testing.T) {
	mockRepo := new(MockStudentRepo)
	service := services.NewStudentService(mockRepo)

	// Case 1: Success
	// Service usually checks existence first
	mockRepo.On("GetByID", uint(1)).Return(&models.Student{}, nil).Once()
	mockRepo.On("Delete", uint(1)).Return(nil).Once()
	err := service.DeleteStudent(1)
	assert.NoError(t, err)

	// Case 2: Repo Error (e.g., Delete fails)
	mockRepo.On("GetByID", uint(2)).Return(&models.Student{}, nil).Once()
	mockRepo.On("Delete", uint(2)).Return(errors.New("delete failed")).Once()
	err = service.DeleteStudent(2)
	assert.Error(t, err)
}
