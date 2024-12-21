// cmd/webapi/cors.go
package webapi

import (
	"github.com/gorilla/handlers"
	"net/http"
)

// SetupCORS applies a CORS policy to the router. CORS (Cross-Origin Resource Sharing) is a security
// feature that blocks JavaScript requests going across different domains unless specified in a policy.
// This function sends the CORS policy for the texting application API.
func SetupCORS(handler http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
		}),
		handlers.AllowedMethods([]string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		}),
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins for now. Update as needed for security.
		handlers.MaxAge(86400), // Cache the CORS preflight request for 1 day (86400 seconds).
	)(handler)
}
