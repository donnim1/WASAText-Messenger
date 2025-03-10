package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	UserID   string `json:"id"`                 // User identifier (typically derived from auth, here passed explicitly)
	PhotoUrl string `json:"photoUrl,omitempty"` // New photo URL to update to
}

// setMyPhoto updates the user's profile photo.
// ✅ Fix the `setMyPhoto` function to handle both URL and file upload.
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Extract authenticated user ID from JWT or session
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Try handling file upload
	file, header, err := r.FormFile("photo")
	if err == nil { // No error → File upload is happening
		defer file.Close() // Ensure the file is closed properly

		// Ensure upload directory exists
		uploadDir := "uploads/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			log.Printf("❌ Failed to create upload directory: %v", err)
			http.Error(w, "Server error: cannot create directory", http.StatusInternalServerError)
			return
		}

		// Save file to uploads directory
		filePath := filepath.Join(uploadDir, header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			log.Printf("❌ Failed to create file: %v", err)
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// Copy file contents
		if _, err := io.Copy(out, file); err != nil {
			log.Printf("❌ Failed to write file: %v", err)
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		// Save file URL to database
		photoUrl := fmt.Sprintf("/%s", filePath)
		if err := rt.db.UpdateUserPhoto(userID, photoUrl); err != nil {
			log.Printf("❌ Database update failed: %v", err)
			http.Error(w, "Failed to update photo in database", http.StatusInternalServerError)
			return
		}

		// Respond with the new photo URL
		response := map[string]string{"photoUrl": photoUrl}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}

		log.Println("✅ Photo successfully updated:", photoUrl)
		return
	}

	// If file upload fails, try JSON-based photo update
	var req userPhotoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
		if req.PhotoUrl == "" {
			log.Println("❌ Invalid photo URL received")
			http.Error(w, "Invalid photo URL", http.StatusBadRequest)
			return
		}

		// Update database
		if err := rt.db.UpdateUserPhoto(userID, req.PhotoUrl); err != nil {
			log.Printf("❌ Database update failed: %v", err)
			http.Error(w, "Failed to update photo URL", http.StatusInternalServerError)
			return
		}

		// Respond with success
		response := map[string]string{"photoUrl": req.PhotoUrl}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}

		log.Println("✅ Photo URL successfully updated:", req.PhotoUrl)
		return
	}

	log.Printf("❌ JSON decoding error: %v", err)
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

// listUsersResponse defines the JSON response for listing users.
type listUsersResponse struct {
	Users []UserSummary `json:"users"`
}

// UserSummary represents a simplified user object.
type UserSummary struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoUrl string `json:"photoUrl"`
}

// listUsers handles GET requests to /users and returns all users.
func (rt *_router) listUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// (Optional) Validate Authorization header if you want only authenticated users to view the list.
	_, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the database function to list users.
	users, err := rt.db.ListUsers()
	if err != nil {
		http.Error(w, "Failed to list users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build a response slice with just the summary info.
	var summaries []UserSummary
	for _, u := range users {
		photo := ""
		if u.PhotoUrl.Valid {
			photo = u.PhotoUrl.String
		}
		summaries = append(summaries, UserSummary{
			ID:       u.ID,
			Username: u.Username,
			PhotoUrl: photo,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listUsersResponse{Users: summaries}); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
