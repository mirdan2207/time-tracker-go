package migrations

import (
	"log"
	// "time"
	"time-tracker-go/models"

	// "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Migrate performs database schema migration for User, Task, and People models.
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.People{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}
