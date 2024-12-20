package routers

import (
	"net/http"
	"github.com/WASAText/pkg/controllers"

	"github.com/gorilla/mux"
)

// InitRouter initializes the router and registers routes
func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// Register login endpoint
	router.HandleFunc("/session", controllers.LoginHandler).Methods(http.MethodPost)

	return router
}
3
