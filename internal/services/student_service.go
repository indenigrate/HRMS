package services

import (
	"hrms_backend/internal/models"
	"hrms_backend/internal/repository"
	"hrms_backend/internal/viewmodels"
)

type StudentService interface {
	CreateStudent(req viewmodels.CreateStudentRequest) (*viewmodels.StudentResponse, error)
	GetAllStudents() ([]viewmodels.StudentResponse, error)
}

type studentService struct {
	repo repository.StudentRepository
}

// Constructor
func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

// Step E: The Logic Layer
func (s *studentService) CreateStudent(req viewmodels.CreateStudentRequest) (*viewmodels.StudentResponse, error) {

	// 1️⃣ Convert ViewModel → Model
	student := models.Student{
		Name:       req.Name,
		Email:      req.Email,
		Department: req.Department,
	}

	// 2️⃣ Call Repository
	err := s.repo.Create(&student)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Convert Model → Response DTO
	response := viewmodels.StudentResponse{
		ID:         student.ID,
		Name:       student.Name,
		Email:      student.Email,
		Department: student.Department,
		CreatedAt:  student.CreatedAt,
	}

	return &response, nil
}

// Get all students
func (s *studentService) GetAllStudents() ([]viewmodels.StudentResponse, error) {
	students, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	responses := make([]viewmodels.StudentResponse, 0, len(students))
	for _, st := range students {
		responses = append(responses, viewmodels.StudentResponse{
			ID:         st.ID,
			Name:       st.Name,
			Email:      st.Email,
			Department: st.Department,
			CreatedAt:  st.CreatedAt,
		})
	}

	return responses, nil
}
