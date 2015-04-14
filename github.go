package edb

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
)

// GitHubFetcher periodically fetches new events for a user.
type GitHubFetcher struct {
	Client   *github.Client
	DB       *DB
	Username string

	Logger *log.Logger
}

// NewGitHubFetcher returns a new GitHubFetcher for a user.
func NewGitHubFetcher(client *github.Client, db *DB, username string) *GitHubFetcher {
	return &GitHubFetcher{
		Client:   client,
		DB:       db,
		Username: username,

		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (f *GitHubFetcher) Run(closing chan struct{}) {
	start := make(chan struct{}, 1)
	start <- struct{}{}

	ticker := time.NewTicker(60 * time.Second)
	for {
		// Wait for next tick or close signal.
		select {
		case <-closing:
			return
		case <-start:
		case <-ticker.C:
		}

		f.Logger.Printf("fetching %s", f.Username)

		// Fetch events from github.
		a, resp, err := f.Client.Activity.ListEventsPerformedByUser(f.Username, true, nil)
		if err != nil {
			f.Logger.Printf("list events: %s", err)
			continue
		} else if resp.StatusCode != http.StatusOK {
			f.Logger.Printf("list events: status=%d", resp.StatusCode)
			continue
		}
		f.Logger.Printf("received %d events (remaining=%d)", len(a), resp.Rate.Remaining)

		// Convert events.
		var events []Event
		for _, ghe := range a {
			e := Event{
				ID:        *ghe.ID,
				Type:      *ghe.Type,
				Timestamp: *ghe.CreatedAt,
				Username:  f.Username,
			}
			if ghe.Actor != nil {
				e.Actor = *ghe.Actor.Login
			}
			if ghe.Repo != nil {
				e.Repository = *ghe.Repo.Name
			}

			events = append(events, e)
		}

		// Save events to database.
		if err := f.DB.SaveEvents(events); err != nil {
			f.Logger.Printf("save events: %s", err)
			continue
		}
	}
}
