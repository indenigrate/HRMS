package controllers

import (
	"hrms_backend/internal/services"
	"hrms_backend/internal/viewmodels"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTP for students.
type StudentController struct {
	service services.StudentService
}

// Constructor
func NewStudentController(svc services.StudentService) *StudentController {
	return &StudentController{service: svc}
}

// Register routes under a router group (e.g., /students)
func (ctl *StudentController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", ctl.CreateStudent)
	rg.GET("", ctl.GetAllStudents)
	rg.GET("/:id", ctl.GetStudentByID)
	rg.PUT("/:id", ctl.UpdateStudent)
	rg.DELETE("/:id", ctl.DeleteStudent)
}

// POST /students
func (ctl *StudentController) CreateStudent(c *gin.Context) {
	var req viewmodels.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ctl.service.CreateStudent(req)
	if err != nil {
		// service returned an error (e.g., DB error, validation error)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GET /students
func (ctl *StudentController) GetAllStudents(c *gin.Context) {
	responses, err := ctl.service.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses)
}

// GET /students/:id
func (ctl *StudentController) GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	student, err := ctl.service.GetStudentByID(uint(id))
	if err != nil {
		// Distinguish not-found vs other errors in service if desired.
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

// PUT /students/:id
func (ctl *StudentController) UpdateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req viewmodels.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update. Service should return updated DTO or error.
	updated, err := ctl.service.UpdateStudent(uint(id), req)
	if err != nil {
		// service may return not-found or validation/db error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteStudent handles DELETE /students/:id
func (ctl *StudentController) DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := ctl.service.DeleteStudent(uint(id)); err != nil {
		// service error (not found, DB error, etc.)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 204 No Content is common for successful delete with no body
	c.Status(http.StatusNoContent)
}
