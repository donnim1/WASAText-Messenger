package api

import (
	"errors"
	"net/http"
)

// ErrUnauthorized is returned when the authorization header is missing or invalid.
var ErrUnauthorized = errors.New("unauthorized")

// getAuthenticatedUserID extracts the user ID from the Authorization header.
// It expects the header to be in the format "Bearer <userID>".
func (rt *_router) getAuthenticatedUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	// Check that the header is long enough and starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", ErrUnauthorized
	}
	// Extract the token (here, just the userID)
	userID := authHeader[7:]
	if userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}
