package github

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benbjohnson/edb"
)

// Client connects to GitHub.
type Client struct{}

// Events returns a list of events from GitHub.
func (c *Client) Events() ([]edb.Event, error) {
	// Retrieve GitHub event stream.
	resp, err := http.Get("https://api.github.com/events")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode events.
	var a []event
	if err := json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return nil, err
	}

	// Convert github events to edb events.
	var events []edb.Event
	for _, ghe := range a {
		e := edb.Event{
			ID:       ghe.ID,
			Type:     ghe.Type,
			Username: ghe.Actor.Login,
			Target:   ghe.Repo.URL,
		}

		events = append(events, e)
	}

	return events, nil
}

// event represents an event from GitHub.
type event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Actor     struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
	} `json:"actor"`
	Repo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
}
