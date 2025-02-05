package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// applyCORSHandler applies a CORS policy to the router.
func applyCORSHandler(h http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		handlers.AllowedOrigins([]string{"*"}), // Allow requests from any origin (frontend)
		handlers.AllowCredentials(),
		handlers.MaxAge(1),
	)(h)
}
