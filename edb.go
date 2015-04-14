package edb

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/boltdb/bolt"
)

// DB represents the application level database.
type DB struct {
	db *bolt.DB
}

// NewDB returns a new instance of DB.
func NewDB() *DB {
	return &DB{}
}

// Open opens the underlying database.
func (db *DB) Open(path string) error {
	d, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return err
	}
	db.db = d

	// Initialize top level buckets.
	db.db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte("events"))
		return nil
	})

	return nil
}

// Close closes the underlying database.
func (db *DB) Close() error {
	if db.db != nil {
		db.db.Close()
	}
	return nil
}

// SaveEvents stores events in the database.
func (db *DB) SaveEvents(events []Event) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		// Loop over events and insert.
		for _, e := range events {
			// Convert event back to JSON.
			b, err := json.Marshal(&e)
			if err != nil {
				return err
			}

			// Create actor bucket.
			eventsBkt := tx.Bucket([]byte("events"))
			bkt, err := eventsBkt.CreateBucketIfNotExists([]byte(e.Actor))
			if err != nil {
				return err
			}

			// Insert event into database.
			if err := bkt.Put([]byte(e.ID), b); err != nil {
				return err
			}
		}

		return nil
	})
}

// EventsByActor returns a list of events from the database for a single actor.
func (db *DB) EventsByActor(actor string) ([]Event, error) {
	var events []Event
	err := db.db.View(func(tx *bolt.Tx) error {
		// Retrieve actor bucket.
		bkt := tx.Bucket([]byte("events")).Bucket([]byte(actor))
		if bkt == nil {
			return nil
		}

		// Iterate over all events.
		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e Event
			if err := json.Unmarshal(v, &e); err != nil {
				return err
			}
			events = append(events, e)
		}

		return nil
	})

	return events, err
}

// Events returns a list of all events from the database.
func (db *DB) Events() ([]Event, error) {
	var events []Event
	err := db.db.View(func(tx *bolt.Tx) error {
		eventBkt := tx.Bucket([]byte("events"))

		return eventBkt.ForEach(func(actor, _ []byte) error {
			// Retrieve actor bucket.
			bkt := eventBkt.Bucket(actor)

			// Iterate over all events.
			c := bkt.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var e Event
				if err := json.Unmarshal(v, &e); err != nil {
					return err
				}
				events = append(events, e)
			}

			return nil
		})
	})

	// Sort events by timestamp.
	sort.Sort(Events(events))

	return events, err
}

// Event returns a genericized event.
type Event struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	Username   string    `json:"username"`
	Actor      string    `json:"actor"`
	Repository string    `json:"repository"`
}

type Events []Event

func (a Events) Len() int {
	return len(a)
}

func (a Events) Less(i, j int) bool {
	return a[i].Timestamp.Before(a[j].Timestamp)
}

func (a Events) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
