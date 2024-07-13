package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time-tracker-go/config"
	"time-tracker-go/models"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @Summary Get information about a person by passport series and number
// @Description Retrieves information about a person from the database based on passport series and number
// @Tags people
// @Accept json
// @Produce json
// @Param passportSeries query int true "Passport series"
// @Param passportNumber query int true "Passport number"
// @Success 200 {object} models.People
// @Router /info [get]
func GetPeopleInfo(w http.ResponseWriter, r *http.Request) {
	passportSeries, err := strconv.Atoi(r.FormValue("passportSeries"))
	if err != nil {
		http.Error(w, "Invalid passport series", http.StatusBadRequest)
		log.Printf("Error converting passport series: %v", err)
		return
	}
	passportNumber, err := strconv.Atoi(r.FormValue("passportNumber"))
	if err != nil {
		http.Error(w, "Invalid passport number", http.StatusBadRequest)
		log.Printf("Error converting passport number: %v", err)
		return
	}

	cfg := config.LoadConfig()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Printf("Database connection error: %v", err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Printf("Database connection error: %v", err)
		return
	}
	defer sqlDB.Close()

	var person models.People
	if err := db.Table("peoples").Where("passport_series = ? AND passport_number = ?", passportSeries, passportNumber).First(&person).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Person not found", http.StatusNotFound)
			log.Printf("Person with passport series %d and passport number %d not found", passportSeries, passportNumber)
		} else {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			log.Printf("Database query error: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		log.Printf("JSON encoding error: %v", err)
		return
	}
	log.Printf("Successfully retrieved person info: %+v", person)
}

// @title Time Tracker API
// @version 1.0
// @description API for managing time tracking information
// @BasePath /api/v1
func SetupHandlers(r *mux.Router) {
	r.HandleFunc("/info", GetPeopleInfo).Methods("GET")
}
