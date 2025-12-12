package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"hrms_backend/internal/config"
	"hrms_backend/internal/controllers"
	"hrms_backend/internal/repository"
	"hrms_backend/internal/services"
)

func main() {
	// load env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found (relying on system env vars)")
	}

	// Connect to db
	config.ConnectDB()

	// Initialize the Router
	r := gin.Default()

	// first Repository (Talks to DB)
	// internal/repository/student_repository.go
	studentRepo := repository.NewStudentRepository(config.DB)

	// Service (Talks to Repository)
	// internal/services/student_service.go
	studentService := services.NewStudentService(studentRepo)

	// Controller (Talks to Service)
	// internal/controllers/student_controller.go
	studentController := controllers.NewStudentController(studentService)

	// Create : http://localhost:8080/students
	studentGroup := r.Group("/students")

	// Pass the group to the controller so it can define endpoints
	studentController.RegisterRoutes(studentGroup)

	// 6. Start the server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
