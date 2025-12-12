package models

import "gorm.io/gorm"

// Student maps to the `students` table.
type Student struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100);not null"`
	Email      string `gorm:"type:varchar(150);unique;not null"`
	Department string `gorm:"type:varchar(100)"`
}
