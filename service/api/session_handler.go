package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// loginRequest defines the expected request body.
type loginRequest struct {
	Username string `json:"name"`
}

// loginResponse defines the response structure.
type loginResponse struct {
	Identifier string `json:"identifier"`
}

func (rt *_router) postSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Parse the request
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2. Validate the username length
	if len(req.Username) < 3 || len(req.Username) > 16 {
		http.Error(w, "Username must be 3-16 characters", http.StatusBadRequest)
		return
	}

	// 3. Check if user already exists in the database.
	user, err := rt.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var userID string
	if user != nil {
		// User exists. Use existing identifier.
		userID = user.ID
	} else {
		// If user does not exist, create a new user.
		userID, err = rt.db.CreateUser(req.Username)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	// 4. Return the user identifier.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(loginResponse{Identifier: userID}); err != nil {
		log.Printf("Error encoding login response: %v", err)
	}
}
