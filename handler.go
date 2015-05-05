package edb

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/benbjohnson/edb/assets"
)

//go:generate go-bindata -ignore bindata.go -o assets/bindata.go -pkg assets -prefix assets assets

// Handler represents the HTTP interface to the database.
type Handler struct {
	DB *DB

	// Enables the use of the local file system for assets, when true.
	LocalMode bool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/assets/") {
		h.serveAsset(w, r, strings.TrimPrefix(r.URL.Path, "/assets/"))
		return
	}

	switch r.URL.Path {
	case "/":
		h.serveAsset(w, r, "index.html")
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

// Serves a file from the local file system or embedded in the binary.
func (h *Handler) serveAsset(w http.ResponseWriter, r *http.Request, filename string) {
	// Serve from local file system in local mode.
	if h.LocalMode {
		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
		http.ServeFile(w, r, filepath.Join("assets", filename))
		return
	}

	// Otherwise serve from embedded assets.
	b, err := assets.Asset(filename)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set content type and write file.
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	if _, err := w.Write(b); err != nil {
		log.Printf("serve asset: %s", err)
		return
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
