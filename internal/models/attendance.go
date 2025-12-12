package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	// Foreign Key
	StudentID uint `gorm:"index"` // Index for FK lookups

	Student Student `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Date   time.Time `gorm:"not null"`
	Status string    `gorm:"type:varchar(20);default:'present'"`
}
