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

// / Update the loginResponse struct to include username and photo URL.
type loginResponse struct {
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
	PhotoURL   string `json:"photoUrl,omitempty"`
}

func (rt *_router) postSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user exists.
	user, err := rt.db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var userID, username, photoURL string

	if user != nil {
		// Existing user: retrieve stored details.
		userID = user.ID
		username = user.Username
		if user.PhotoUrl.Valid {
			photoURL = user.PhotoUrl.String
		}
	} else {
		// Create new user.
		userID, err = rt.db.CreateUser(req.Username)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		username = req.Username // New users use the provided username.
		photoURL = ""           // New users have no photo initially.
	}

	// Return the login response with user details.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(loginResponse{
		Identifier: userID,
		Username:   username,
		PhotoURL:   photoURL,
	}); err != nil {
		log.Printf("Error encoding login response: %v", err)
	}
}
