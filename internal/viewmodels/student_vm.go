package viewmodels

import "time"

// when creating a student.
// `binding:"required"` tags enforce presence of the fields.
type CreateStudentRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Department string `json:"department" binding:"required"`
}

// for PUT /students/:id.
type UpdateStudentRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

// GET/POST responses
type StudentResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}
