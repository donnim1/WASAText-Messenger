package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Username string `json:"name"`
}

type loginResponse struct {
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
	PhotoURL   string `json:"photoUrl,omitempty"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use errors.Is to check for wrapped errors (e.g., io.EOF)
		if errors.Is(err, io.EOF) {
			http.Error(w, "Empty request body", http.StatusBadRequest)
		} else {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}
		return
	}

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
		username = req.Username
		photoURL = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Check the error from encoding the response; if an error occurs log it.
	if err := json.NewEncoder(w).Encode(loginResponse{
		Identifier: userID,
		Username:   username,
		PhotoURL:   photoURL,
	}); err != nil {
		log.Printf("Error encoding login response: %v", err)
		// Optionally: you could also call http.Error here, but keep in mind headers are already written.
	}
}
