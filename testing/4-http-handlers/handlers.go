// Package handlers holds small HTTP handlers to demonstrate testing web code
// with net/http/httptest — no real network or ports required.
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Hello writes a fixed greeting.
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// Echo writes the request body back to the client.
func Echo(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(body)
}

// Greet reads a "name" query parameter and responds with JSON. If "name" is
// missing it returns 400.
func Greet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello, " + name + "!",
	})
}
