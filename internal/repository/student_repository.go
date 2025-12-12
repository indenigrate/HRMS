package repository

import (
	"hrms_backend/internal/models"

	"gorm.io/gorm"
)

// handles DB operations (Create, Read, etc.).
// controller depends on this abstraction, not the implementation.

type StudentRepository interface {
	Create(student *models.Student) error
	GetAll() ([]models.Student, error)
	Update(id uint, student *models.Student) error
	GetByID(id uint) (*models.Student, error)
	Delete(id uint) error
}

// the interface
type studentRepo struct {
	db *gorm.DB
}

// Constructor
func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepo{db: db}
}

// Create a student
func (r *studentRepo) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

// Get all students
func (r *studentRepo) GetAll() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Find(&students).Error
	return students, err
}

// 6Get a student by ID (useful for update/delete)
func (r *studentRepo) GetByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// Update a student
func (r *studentRepo) Update(id uint, student *models.Student) error {
	return r.db.Model(&models.Student{}).Where("id = ?", id).Updates(student).Error
}

// Delete a student
func (r *studentRepo) Delete(id uint) error {
	return r.db.Delete(&models.Student{}, id).Error
}
