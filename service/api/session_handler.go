package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Login request format
type loginRequest struct {
	Username string `json:"name"`
}

// Login response format
type loginResponse struct {
	Identifier string `json:"identifier"`
}

func (rt *_router) postSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Parse request
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2. Validate username
	if len(req.Username) < 3 || len(req.Username) > 16 {
		http.Error(w, "Username must be 3-16 characters", http.StatusBadRequest)
		return
	}

	// 3. Check if user exists
	user, err := rt.db.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user != nil {

		// return user ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(loginResponse{
			Identifier: user.ID,
		})
		return
	}

	// 5. If user does not exist, create a new user
	userID, err := rt.db.CreateUser(req.Username)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// 6. Return response with newly created user ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loginResponse{
		Identifier: userID,
	})
}
