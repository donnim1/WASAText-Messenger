package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"WASA/WASAText/models"

	"github.com/gorilla/mux"
)

// GetUsers retrieves a list of users with optional search and limit parameters
func GetUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limitParam := r.URL.Query().Get("limit")

	// Default limit to 20 if not specified
	limit := 20
	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 1 || limit > 100 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	users, err := models.GetAllUsers(search, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve users: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUserName updates the username of the current user
func UpdateUserName(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewName string `json:"newName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := mux.Vars(r)["id"]
	if err := models.UpdateUserName(userID, req.NewName); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update username: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Username updated successfully."})
}

// UpdateUserPhoto updates the user's profile photo
func UpdateUserPhoto(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	// Simulate handling a photo upload and generate a photo URL
	// (In real-world apps, you would save the photo to a file storage service)
	photoURL := "https://example.com/photo/" + userID

	if err := models.UpdateUserPhoto(userID, photoURL); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update photo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Profile photo updated successfully.",
		"photoUrl": photoURL,
	})
}
