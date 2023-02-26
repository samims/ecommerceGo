package utils

import (
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := ToJSON(payload, w)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "internal server error")
	}

}
