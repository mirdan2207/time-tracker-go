package models

import (
	"gorm.io/gorm"
	"time"
)

// Task represents a task assigned to a user.
type Task struct {
	gorm.Model            // Default GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt)
	UserID      uint      `json:"userID"`      // ID of the user associated with the task
	Description string    `json:"description"` // Description of the task
	StartTime   time.Time `json:"startTime"`   // Start time of the task
	EndTime     time.Time `json:"endTime"`     // End time of the task
	Duration    int       `json:"duration"`    // Duration of the task in minutes
}
