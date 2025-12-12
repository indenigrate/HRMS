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
	rg.POST("", ctl.MarkAttendance)
	rg.GET("/student/:studentID", ctl.GetAttendanceByStudentID)
}

// POST /attendance
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

// GET /attendance/student/:studentID
func (ctl *AttendanceController) GetAttendanceByStudentID(c *gin.Context) {
	idStr := c.Param("studentID")
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
