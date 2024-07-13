package models

import "gorm.io/gorm"

// People represents a person with personal details.
type People struct {
	gorm.Model            // Default GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt)
	PassportSeries int    `gorm:"unique;not null" json:"passportSeries"` // Passport series (unique, not null)
	PassportNumber int    `gorm:"unique;not null" json:"passportNumber"` // Passport number (unique, not null)
	Surname        string `json:"surname"`                               // Surname of the person
	Name           string `json:"name"`                                  // Name of the person
	Patronymic     string `json:"patronymic"`                            // Patronymic (middle name) of the person
	Address        string `json:"address"`                               // Address of the person
}
