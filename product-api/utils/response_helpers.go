package utils

import (
	"net/http"
)

// RespondWithError is a helper function that writes a JSON-encoded
// error response to the HTTP response writer. It sets the status code
// and content type headers and encodes the provided error message in
// a GenericError struct.

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonErr := ToJSON(GenericError{Message: message}, w)
	if jsonErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
