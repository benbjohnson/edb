package edb

import (
	"encoding/json"

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
	return db.db.Close()
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

			// Insert event into database.
			bkt := tx.Bucket([]byte("events"))
			if err := bkt.Put([]byte(e.ID), b); err != nil {
				return err
			}
		}

		return nil
	})
}

// Events returns a list of events from the database.
func (db *DB) Events() ([]Event, error) {
	var events []Event
	err := db.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("events"))
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

// Event returns a genericized event.
type Event struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
	Target   string `json:"target"`
}
