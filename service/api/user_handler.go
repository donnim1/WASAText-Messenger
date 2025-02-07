package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// userNameUpdateRequest defines the expected JSON payload for updating the username.
// We no longer require the user ID in the request because it is obtained from the Authorization header.
type userNameUpdateRequest struct {
	NewName string `json:"newName"` // New username to update to
}

// userNameUpdateResponse defines the response after updating the username.
type userNameUpdateResponse struct {
	Message string `json:"message"`
}

// setMyUserName updates the username for an existing user.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Validate the Authorization header and extract the authenticated user ID.
	authenticatedUserID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body.
	var req userNameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate the new username length.
	if len(req.NewName) < 3 || len(req.NewName) > 16 {
		http.Error(w, "Username must be 3-16 characters", http.StatusBadRequest)
		return
	}

	// Update the username in the database using the authenticated user ID.
	err = rt.db.UpdateUserName(authenticatedUserID, req.NewName)
	if err != nil {
		http.Error(w, "Failed to update username: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userNameUpdateResponse{
		Message: "Username updated successfully",
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// userPhotoUpdateRequest defines the expected JSON payload for updating the user photo.
type userPhotoUpdateRequest struct {
	UserID   string `json:"id"`       // User identifier (typically derived from auth, here passed explicitly)
	PhotoURL string `json:"photoUrl"` // New photo URL to update to
}

// userPhotoUpdateResponse defines the response after updating the profile photo.
type userPhotoUpdateResponse struct {
	Message string `json:"message"`
}

// setMyPhoto updates the user's profile photo.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Validate the Authorization header
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// 1. Parse the request body
	var req userPhotoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 2. Validate that PhotoURL is not empty
	if req.PhotoURL == "" {
		http.Error(w, "Photo URL cannot be empty", http.StatusBadRequest)
		return
	}

	// 3. Update the photo in the database
	err = rt.db.UpdateUserPhoto(userID, req.PhotoURL)
	if err != nil {
		http.Error(w, "Failed to update photo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userPhotoUpdateResponse{
		Message: "Photo updated successfully",
	}); err != nil {
		// Log the error. The response is already sent, so this is only for debugging.
		log.Printf("Error encoding response: %v", err)
	}

}
