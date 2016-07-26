package restutil

import (
	"encoding/json"
	"net/http"
)

// Writes body as json to the response writer
func WriteJson(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(body)
}

// Reads json body from the request
func ReadJson(r *http.Request, v interface{}) {
	// TODO deal with an error
	json.NewDecoder(r.Body).Decode(v)
}
