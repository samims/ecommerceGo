package utils

import (
	"net/http"
)

// RespondWithJSON sends a JSON-encoded response to the client with the specified HTTP status code and payload.
// Arguments:
//
//	w:  `http.ResponseWriter` is writer of the response
//
//	code: integer http status code for the resp
//
//	payload: parameter is an interface{} type that represents the data written to body
//
// If the encoding of the `payload` fails, this function sends an error response with
// HTTP status code 500 and a message "internal server error".
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := ToJSON(payload, w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal server error")
	}

}
