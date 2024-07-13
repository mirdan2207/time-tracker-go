package models

import "gorm.io/gorm"

// User represents a user in the system.
type User struct {
	gorm.Model            // Default GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt)
	PassportNumber string `gorm:"unique;not null" json:"passportNumber"` // Passport number of the user (unique and not null)
	Surname        string `json:"surname"`                               // Surname of the user
	Name           string `json:"name"`                                  // Name of the user
	Patronymic     string `json:"patronymic"`                            // Patronymic (middle name) of the user
	Address        string `json:"address"`                               // Address of the user
	Tasks          []Task `json:"tasks"`                                 // List of tasks associated with the user
}
