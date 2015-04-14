package edb

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler represents the HTTP interface to the database.
type Handler struct {
	DB *DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/events.json":
		if r.Method == "GET" {
			h.serveEvents(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}

// serveEvents writes all events in the database to the handler.
func (h *Handler) serveEvents(w http.ResponseWriter, r *http.Request) {
	// Retrieve list of events.
	a, err := h.DB.Events()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write events to response.
	if err := json.NewEncoder(w).Encode(a); err != nil {
		log.Print(err)
	}
}
