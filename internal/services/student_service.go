package services

import (
	"hrms_backend/internal/models"
	"hrms_backend/internal/repository"
	"hrms_backend/internal/viewmodels"
)

type StudentService interface {
	CreateStudent(req viewmodels.CreateStudentRequest) (*viewmodels.StudentResponse, error)
	GetAllStudents() ([]viewmodels.StudentResponse, error)
	GetStudentByID(id uint) (*viewmodels.StudentResponse, error)
	UpdateStudent(id uint, req viewmodels.UpdateStudentRequest) (*viewmodels.StudentResponse, error)
	DeleteStudent(id uint) error
}

type studentService struct {
	repo repository.StudentRepository
}

// Constructor
func NewStudentService(repo repository.StudentRepository) StudentService {
	// sends
	return &studentService{repo: repo}
}

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

func (s *studentService) GetStudentByID(id uint) (*viewmodels.StudentResponse, error) {
	st, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	resp := viewmodels.StudentResponse{
		ID:         st.ID,
		Name:       st.Name,
		Email:      st.Email,
		Department: st.Department,
		CreatedAt:  st.CreatedAt,
	}
	return &resp, nil
}

// UpdateStudent updates fields provided in the request and returns the updated DTO.
// It reads the existing record, updates only non-empty fields
func (s *studentService) UpdateStudent(id uint, req viewmodels.UpdateStudentRequest) (*viewmodels.StudentResponse, error) {
	// fetch existing
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// apply changes only when provided (empty string => no change)
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Department != "" {
		existing.Department = req.Department
	}

	// persist update
	if err := s.repo.Update(id, existing); err != nil {
		return nil, err
	}

	// fetch again to ensure fields like UpdatedAt are current (optional)
	updated, err := s.repo.GetByID(id)
	if err != nil {
		// if fetch fails after update, surface the update error or return a wrapped error
		return nil, err
	}

	resp := viewmodels.StudentResponse{
		ID:         updated.ID,
		Name:       updated.Name,
		Email:      updated.Email,
		Department: updated.Department,
		CreatedAt:  updated.CreatedAt,
	}
	return &resp, nil
}

// deletes the student with the given id.
func (s *studentService) DeleteStudent(id uint) error {
	// Optionally, verify existence first
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
