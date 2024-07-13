package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time-tracker-go/config"
	"time-tracker-go/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserController handles HTTP requests related to users.
type UserController struct {
	DB     *gorm.DB
	Config config.Config
}

// NewUserController creates a new instance of UserController with the given DB connection and configuration.
func NewUserController(db *gorm.DB, config config.Config) *UserController {
	return &UserController{DB: db, Config: config}
}

type AddUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}

// @Summary Get users with optional filters and pagination
// @Description Retrieves users based on optional filters and supports pagination
// @Tags users
// @Accept json
// @Produce json
// @Param passportNumber query string false "Passport number"
// @Param surname query string false "Surname"
// @Param name query string false "Name"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {array} models.User
// @Router /users [get]
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	query := uc.DB

	// Filtration
	if passportNumber := r.URL.Query().Get("passportNumber"); passportNumber != "" {
		query = query.Where("passport_number = ?", passportNumber)
	}

	if surname := r.URL.Query().Get("surname"); surname != "" {
		query = query.Where("surname = ?", surname)
	}

	if name := r.URL.Query().Get("name"); name != "" {
		query = query.Where("name = ?", name)
	}

	if patronymic := r.URL.Query().Get("patronymic"); patronymic != "" {
		query = query.Where("patronymic = ?", patronymic)
	}

	if address := r.URL.Query().Get("address"); address != "" {
		query = query.Where("address = ?", address)
	}

	// Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query = query.Limit(pageSize).Offset(offset).Find(&users)

	if query.Error != nil {
		http.Error(w, query.Error.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching users: %v", query.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

	log.Printf("Fetched %d users", len(users))
}

// @Summary Delete a user by ID
// @Description Deletes a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("User not found with ID %d", id)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching user: %v", err)
		return
	}

	if err := uc.DB.Delete(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error deleting user: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})

	log.Printf("Deleted user with ID %d", id)
}

// @Summary Update a user by ID
// @Description Updates a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user object"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Printf("User not found with ID %d", id)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching user: %v", err)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Invalid request payload: %v", err)
		return
	}

	user.PassportNumber = updatedUser.PassportNumber
	user.Surname = updatedUser.Surname
	user.Name = updatedUser.Name
	user.Patronymic = updatedUser.Patronymic
	user.Address = updatedUser.Address

	if err := uc.DB.Save(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error updating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	log.Printf("Updated user with ID %d", id)
}

// @Summary Add a new user
// @Description Adds a new user based on the provided passport number, fetching details from an external API
// @Tags users
// @Accept json
// @Produce json
// @Param request body AddUserRequest true "Request body with passport number"
// @Success 201 {object} models.User
// @Router /users [post]
func (uc *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	var request AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Invalid request payload: %v", err)
		return
	}

	parts := strings.Split(request.PassportNumber, " ")
	if len(parts) != 2 {
		http.Error(w, "Invalid passport number format", http.StatusBadRequest)
		return
	}

	passportSeries, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid passport series", http.StatusBadRequest)
		log.Printf("Invalid passport series: %v", err)
		return
	}

	passportNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid passport number", http.StatusBadRequest)
		log.Printf("Invalid passport number: %v", err)
		return
	}

	apiURL := fmt.Sprintf(uc.Config.ExternalAPIURL+"/info?passportSeries=%d&passportNumber=%d", passportSeries, passportNumber)
	log.Printf("Calling external API: %s", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch user info from external API", http.StatusInternalServerError)
		log.Printf("Failed to fetch user info from external API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch user info from external API", http.StatusInternalServerError)
		log.Printf("Failed to fetch user info from external API. Status code: %d", resp.StatusCode)
		return
	}

	var people models.People
	if err := json.NewDecoder(resp.Body).Decode(&people); err != nil {
		http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
		log.Printf("Failed to decode API response: %v", err)
		return
	}

	user := models.User{
		PassportNumber: request.PassportNumber,
		Surname:        people.Surname,
		Name:           people.Name,
		Patronymic:     people.Patronymic,
		Address:        people.Address,
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	log.Printf("Created user with Passport Number %s", user.PassportNumber)
}
