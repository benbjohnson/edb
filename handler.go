package edb

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Handler represents the HTTP interface to the database.
type Handler struct {
	DB *DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/assets/") {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "assets/index.html")
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

	// Marshal events with pretty printing.
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Allow any front end to access our data.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write JSON out to response.
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Print(err)
	}
}
