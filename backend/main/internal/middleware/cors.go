package middleware

import (
	"net/http"
)

// CORS middleware function to add CORS headers
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Change to specific origin if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the method is OPTIONS, return without calling the next handler
		if r.Method == http.MethodOptions {
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
