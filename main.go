package main

import (
	"hrms_backend/internal/config"
	"hrms_backend/internal/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()

	// auto create tables if not exists
	config.DB.AutoMigrate(&models.Student{})

}
