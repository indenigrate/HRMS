package repository

import (
	"hrms_backend/internal/models"

	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	GetAttendanceByStudentID(studentID uint) ([]models.Attendance, error)
}

type attendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepo{db: db}
}

func (r *attendanceRepo) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

// GetAttendanceByStudentID fetches attendance and Preloads the Student details
func (r *attendanceRepo) GetAttendanceByStudentID(studentID uint) ([]models.Attendance, error) {
	var attendanceList []models.Attendance

	// Uses Preload to fetch the associated Student entity in an optimized way
	err := r.db.Preload("Student").Where("student_id = ?", studentID).Find(&attendanceList).Error

	return attendanceList, err
}
