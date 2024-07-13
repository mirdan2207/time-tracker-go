package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"time-tracker-go/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// TaskController handles HTTP requests related to tasks.
type TaskController struct {
	DB *gorm.DB
}

// NewTaskController creates a new instance of TaskController with the given DB connection.
func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{DB: db}
}

// @Summary Get time entries by user ID and period
// @Description Retrieves time entries for a user within a specified time period
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param start_date query string true "Start date (format: 2006-01-02T15:04:05)"
// @Param end_date query string true "End date (format: 2006-01-02T15:04:05)"
// @Success 200 {array} models.Task
// @Router /users/{id}/tasks [get]
func (tc *TaskController) GetTimeEntriesByUserAndPeriod(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing start_date and end_date", http.StatusBadRequest)
		log.Println("Missing start_date or end_date")
		return
	}

	startDate, err := time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		log.Printf("Invalid start_date format: %v", err)
		return
	}

	endDate, err := time.Parse("2006-01-02T15:04:05", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end_date format", http.StatusBadRequest)
		log.Printf("Invalid end_date format: %v", err)
		return
	}

	log.Printf("Fetching time entries for user %d between %s and %s", userID, startDateStr, endDateStr)

	var tasks []models.Task
	if err := tc.DB.Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).Order("duration DESC").Find(&tasks).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching time entries: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

	log.Printf("Successfully fetched %d time entries", len(tasks))
}

// @Summary Start a task for a user
// @Description Starts a task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param taskID path int true "Task ID"
// @Success 200 {object} models.Task
// @Router /users/{id}/tasks/{taskID}/start [post]
func (tc *TaskController) StartTaskForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	taskID, err := strconv.Atoi(params["taskID"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Printf("Invalid task ID: %v", err)
		return
	}

	var task models.Task
	if err := tc.DB.Where("user_id=? AND id=?", userID, taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			log.Printf("Task not found for user %d with task ID %d", userID, taskID)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching task: %v", err)
		return
	}

	task.StartTime = time.Now()

	if err := tc.DB.Save(&task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error saving task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)

	log.Printf("Task %d started for user %d", taskID, userID)
}

// @Summary End a task for a user
// @Description Ends a task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param taskID path int true "Task ID"
// @Success 200 {object} models.Task
// @Router /users/{id}/tasks/{taskID}/end [post]
func (tc *TaskController) EndTaskForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	taskID, err := strconv.Atoi(params["taskID"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		log.Printf("Invalid task ID: %v", err)
		return
	}

	var task models.Task
	if err := tc.DB.Where("user_id=? AND id=?", userID, taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			log.Printf("Task not found for user %d with task ID %d", userID, taskID)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching task: %v", err)
		return
	}

	task.EndTime = time.Now()
	task.Duration = int(task.EndTime.Sub(task.StartTime).Minutes())

	if err := tc.DB.Save(&task).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error saving task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)

	log.Printf("Task %d ended for user %d", taskID, userID)
}

// @Summary Add a task for a user
// @Description Adds a new task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param task body models.Task true "Task object to be added"
// @Success 201 {object} models.Task
// @Router /users/{id}/tasks [post]
func (tc *TaskController) AddTaskForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		log.Printf("Invalid user ID: %v", err)
		return
	}

	var newTask models.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Invalid request payload: %v", err)
		return
	}

	newTask.UserID = uint(userID)

	if err := tc.DB.Create(&newTask).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating task: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

	log.Printf("Task created successfully for user %d with ID %d", userID, newTask.ID)
}
