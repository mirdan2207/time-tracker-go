package routes

import (
	"log"
	"net/http"
	"time-tracker-go/api"
	"time-tracker-go/config"
	"time-tracker-go/controllers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
	_ "time-tracker-go/docs" // Import generated docs
)

// SetupRoutes configures the application routes and returns the router.
// Swagger:Route GET /users getUsers
// Get a list of users.
// Responses:
//   200: usersResponse

// Swagger:Route DELETE /users/{id} deleteUser
// Delete a user by ID.
// Parameters:
//   id path int true "User ID"
// Responses:
//   204: noContentResponse

// Swagger:Route PUT /users/{id} updateUser
// Update a user by ID.
// Parameters:
//   id path int true "User ID"
// Responses:
//   200: userResponse

// Swagger:Route POST /users addUser
// Add a new user.
// Responses:
//   201: userResponse

// Swagger:Route GET /users/{id}/time-entries getTimeEntriesByUserAndPeriod
// Get time entries for a user and period.
// Parameters:
//   id path int true "User ID"
// Responses:
//   200: timeEntriesResponse

// Swagger:Route POST /users/{id}/tasks addTaskForUser
// Add a task for a user.
// Parameters:
//   id path int true "User ID"
// Responses:
//   201: taskResponse

// Swagger:Route PUT /users/{id}/tasks/{taskID}/start startTaskForUser
// Start a task for a user.
// Parameters:
//   id path int true "User ID"
//   taskID path int true "Task ID"
// Responses:
//   200: taskResponse

// Swagger:Route PUT /users/{id}/tasks/{taskID}/end endTaskForUser
// End a task for a user.
// Parameters:
//   id path int true "User ID"
//   taskID path int true "Task ID"
// Responses:
//   200: taskResponse

func SetupRoutes(db *gorm.DB, cfg config.Config) *mux.Router {
	router := mux.NewRouter()

	userController := controllers.NewUserController(db, cfg)
	taskController := controllers.NewTaskController(db)

	// Routes for user management
	router.HandleFunc("/users", logRequest(userController.GetUsers)).Methods("GET")
	router.HandleFunc("/users", logRequest(userController.AddUser)).Methods("POST")
	router.HandleFunc("/users/{id}", logRequest(userController.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/users/{id}", logRequest(userController.UpdateUser)).Methods("PUT")

	// Routes for user task management
	router.HandleFunc("/users/{id}/time-entries", logRequest(taskController.GetTimeEntriesByUserAndPeriod)).Methods("GET")
	router.HandleFunc("/users/{id}/tasks", logRequest(taskController.AddTaskForUser)).Methods("POST")
	router.HandleFunc("/users/{id}/tasks/{taskID}/start", logRequest(taskController.StartTaskForUser)).Methods("PUT")
	router.HandleFunc("/users/{id}/tasks/{taskID}/end", logRequest(taskController.EndTaskForUser)).Methods("PUT")

	// Setting up sub-routes for API
	apiRouter := router.PathPrefix("/api").Subrouter()
	api.SetupHandlers(apiRouter)

	// Swagger route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}

// logRequest returns a wrapper for HTTP handler that logs requests.
func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	}
}
