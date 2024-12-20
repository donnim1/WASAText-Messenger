package routes

import (
	"net/http"
	"WASA/WASAText/pkg/controllers"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes registers routes for user operations
func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/users", controllers.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}/name", controllers.UpdateUserName).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}/photo", controllers.UpdateUserPhoto).Methods(http.MethodPut)
}
