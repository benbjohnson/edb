package main

import (
	"fmt"
	"log"

	"github.com/benbjohnson/edb"
	"github.com/benbjohnson/edb/github"
)

func main() {
	db := edb.NewDB()
	if err := db.Open("data"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// load(db)
	printDB(db)
}

func printDB(db *edb.DB) {
	// Retrieve events from database.
	events, err := db.Events()
	if err != nil {
		log.Fatal(err)
	}

	// Print events.
	for _, e := range events {
		fmt.Printf("> %s: %s\n", e.ID, e.Username)
	}
}

func load(db *edb.DB) {
	// Retrieve events from GitHub.
	var c github.Client
	events, err := c.Events()
	if err != nil {
		log.Fatal(err)
	}

	// Save events to database.
	if err := db.SaveEvents(events); err != nil {
		log.Fatal(err)
	}
}
