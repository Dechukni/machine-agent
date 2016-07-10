package rest

import (
	"encoding/json"
	"net/http"
)

const (
	APPLICATION_JSON = "application/json"
)

// Writes body as json to the response writer
func WriteJson(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	// TODO deal with an error
	json.NewEncoder(w).Encode(body)
}

// Read json body from the request
func ReadJson(r *http.Request, v interface{}) {
	// TODO deal with an error
	json.NewDecoder(r.Body).Decode(v)
}
