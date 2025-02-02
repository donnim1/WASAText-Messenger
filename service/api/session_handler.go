package api

import (
	"encoding/json"
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
	
	// 3. Create/get user
	userID, err := rt.db.CreateUser(req.Username)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	
	// 4. Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loginResponse{
		Identifier: userID,
	})
}