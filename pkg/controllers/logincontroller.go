package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/WASA/WASAText/pkg/models"
)

// LoginRequest represents the request payload for login
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse represents the response payload for login
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// LoginHandler handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate username
	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Username must be between 3 and 16 characters", http.StatusBadRequest)
		return
	}

	// Get or create user
	identifier, err := models.GetOrCreateUser(req.Name)
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	// Respond with the identifier
	response := LoginResponse{Identifier: identifier}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
