package controllers

import (
	"hrms_backend/internal/services"
	"hrms_backend/internal/viewmodels"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	service services.AttendanceService
}

func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{service: service}
}

func (ctl *AttendanceController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/mark", ctl.MarkAttendance)
	rg.GET("/:student_id", ctl.GetAttendanceByStudentID)
}

// MarkAttendance handles POST /attendance/mark
// @Summary      Mark student attendance
// @Description  Marks a student's attendance as 'Present' or 'Absent' for a given date.
// @Tags         Attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      viewmodels.CreateAttendanceRequest  true  "Attendance details"
// @Success      201         "Created"
// @Failure      400         {object}  viewmodels.ErrorResponse
// @Router       /attendance/mark [post]
func (ctl *AttendanceController) MarkAttendance(c *gin.Context) {
	var req viewmodels.CreateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctl.service.MarkAttendance(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetAttendanceByStudentID handles GET /attendance/:student_id
// @Summary      Get attendance by student ID
// @Description  Retrieves all attendance records for a specific student.
// @Tags         Attendance
// @Produce      json
// @Param        student_id  path      int  true  "Student ID"
// @Success      200         {array}   viewmodels.AttendanceResponse
// @Failure      400         {object}  viewmodels.ErrorResponse
// @Failure      500         {object}  viewmodels.ErrorResponse
// @Router       /attendance/{student_id} [get]
func (ctl *AttendanceController) GetAttendanceByStudentID(c *gin.Context) {
	idStr := c.Param("student_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	resp, err := ctl.service.GetAttendanceByStudentID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
