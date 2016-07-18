package restutil

import (
	"encoding/json"
	"net/http"
)

const (
	APPLICATION_JSON = "application/json"
)

// Writes body as json to the response writer
func WriteJson(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	return json.NewEncoder(w).Encode(body)
}

// Reads json body from the request
func ReadJson(r *http.Request, v interface{}) {
	// TODO deal with an error
	json.NewDecoder(r.Body).Decode(v)
}
