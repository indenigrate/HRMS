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

// CreateStudent handles POST /students
// @Summary      Create a new student
// @Description  Creates a new student record in the database.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        student  body      viewmodels.CreateStudentRequest  true  "Student details"
// @Success      201      {object}  viewmodels.StudentResponse
// @Failure      400      {object}  viewmodels.ErrorResponse
// @Router       /students [post]
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

// GetAllStudents handles GET /students
// @Summary      Get all students
// @Description  Retrieves a paginated list of all students.
// @Tags         Students
// @Produce      json
// @Param        page   query     int  false  "Page number for pagination"  minimum(1)
// @Param        limit  query     int  false  "Number of items per page"    minimum(1)
// @Success      200    {array}   viewmodels.StudentResponse
// @Failure      500    {object}  viewmodels.ErrorResponse
// @Router       /students [get]
func (ctl *StudentController) GetAllStudents(c *gin.Context) {
	// Parse Query Params with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	responses, err := ctl.service.GetAllStudents(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses)
}

// GetStudentByID handles GET /students/:id
// @Summary      Get a student by ID
// @Description  Retrieves the details of a single student by their unique ID.
// @Tags         Students
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  viewmodels.StudentResponse
// @Failure      400  {object}  viewmodels.ErrorResponse
// @Failure      404  {object}  viewmodels.ErrorResponse
// @Router       /students/{id} [get]
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

// UpdateStudent handles PUT /students/:id
// @Summary      Update a student
// @Description  Updates an existing student's details by their ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        id       path      int                              true  "Student ID"
// @Param        student  body      viewmodels.UpdateStudentRequest  true  "Updated student details"
// @Success      200      {object}  viewmodels.StudentResponse
// @Failure      400      {object}  viewmodels.ErrorResponse
// @Router       /students/{id} [put]
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
// @Summary      Delete a student
// @Description  Deletes a student from the database by their ID.
// @Tags         Students
// @Produce      json
// @Param        id  path  int  true  "Student ID"
// @Success      204 "No Content"
// @Failure      400 {object} viewmodels.ErrorResponse
// @Router       /students/{id} [delete]
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
