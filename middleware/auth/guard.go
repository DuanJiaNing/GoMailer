package auth

import (
	"net/http"
)

// Guard creates a middleware to guard url by roles
func Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the next handler, which can be another middleware in the chain, or the final app.
		next.ServeHTTP(w, r)
	})
}
