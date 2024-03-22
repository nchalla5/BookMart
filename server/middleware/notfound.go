package middleware

import (
	"encoding/json"
	"net/http"
)

// NotFoundHandler returns a custom 404 response in standard HTTP format.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": "The requested resource was not found.",
	})
}
