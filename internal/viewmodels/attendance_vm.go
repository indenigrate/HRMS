package viewmodels

import "time"

type CreateAttendanceRequest struct {
	StudentID uint      `json:"student_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	// oneof validation ensures only valid statuses are accepted
	Status string `json:"status" binding:"required,oneof=present absent late excused"`
}

type AttendanceResponse struct {
	ID          uint      `json:"id"`
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name,omitempty"` // Optional: filled if Student is preloaded
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}
