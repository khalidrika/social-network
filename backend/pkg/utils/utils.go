package utils

import (
	"net/http"

	uuid "github.com/gofrs/uuid"
)

// Create unique ID
func UniqueId() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func ConfigHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
	return w
}

func ConfigFSHeader(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:5173")
		fs.ServeHTTP(w, r)
	}
}
