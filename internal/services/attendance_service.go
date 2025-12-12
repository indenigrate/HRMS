package services

import (
	"errors"
	"hrms_backend/internal/models"
	"hrms_backend/internal/repository"
	"hrms_backend/internal/viewmodels"
)

type AttendanceService interface {
	MarkAttendance(req viewmodels.CreateAttendanceRequest) error
	GetAttendanceByStudentID(studentID uint) ([]viewmodels.AttendanceResponse, error)
}

type attendanceService struct {
	attRepo     repository.AttendanceRepository
	studentRepo repository.StudentRepository // Dependency injected for Logic Check
}

// Constructor: Requires both repositories
func NewAttendanceService(attRepo repository.AttendanceRepository, studentRepo repository.StudentRepository) AttendanceService {
	return &attendanceService{
		attRepo:     attRepo,
		studentRepo: studentRepo,
	}
}

// MarkAttendance handles the business logic for creating attendance
func (s *attendanceService) MarkAttendance(req viewmodels.CreateAttendanceRequest) error {

	// STEP D: Logic Check - Verify Student Exists
	// We use s.studentRepo.GetByID to ensure we don't mark attendance for a non-existent ID.
	_, err := s.studentRepo.GetByID(req.StudentID)
	if err != nil {
		return errors.New("student not found: cannot mark attendance")
	}

	// Logic Check Passed: Create the Model
	attendance := models.Attendance{
		StudentID: req.StudentID,
		Date:      req.Date,
		Status:    req.Status,
	}

	// Persist
	return s.attRepo.Create(&attendance)
}

// GetAttendanceByStudentID fetches records and maps them to ViewModels
func (s *attendanceService) GetAttendanceByStudentID(studentID uint) ([]viewmodels.AttendanceResponse, error) {
	// 1. Verify student exists (Optional, but good practice)
	if _, err := s.studentRepo.GetByID(studentID); err != nil {
		return nil, errors.New("student not found")
	}

	// 2. Fetch data
	records, err := s.attRepo.GetAttendanceByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	// 3. Convert to DTOs
	responses := make([]viewmodels.AttendanceResponse, 0, len(records))
	for _, rec := range records {
		resp := viewmodels.AttendanceResponse{
			ID:        rec.ID,
			StudentID: rec.StudentID,
			Date:      rec.Date,
			Status:    rec.Status,
		}
		// If the Student relation was preloaded in the repo, we can map the name
		if rec.Student.Name != "" {
			resp.StudentName = rec.Student.Name
		}
		responses = append(responses, resp)
	}

	return responses, nil
}
