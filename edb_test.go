package edb_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/edb"
)

// Ensure events can be added and retrieve from the database by actor name.
func TestDB_EventsByActor(t *testing.T) {
	db := edb.NewDB()
	if err := db.Open(MustTempFile()); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Save events.
	if err := db.SaveEvents([]edb.Event{
		{ID: "1", Type: "PushEvent", Timestamp: MustParseTime("2000-01-01T00:00:00Z"), Actor: "bob"},
		{ID: "2", Type: "IssueEvent", Timestamp: MustParseTime("2000-02-01T00:00:00Z"), Actor: "bob"},
		{ID: "3", Type: "PullEvent", Timestamp: MustParseTime("2000-01-06T00:00:00Z"), Actor: "susy"},
	}); err != nil {
		t.Fatalf("save events: %s", err)
	}

	// Retrieve events.
	a, err := db.EventsByActor("bob")
	if err != nil {
		t.Fatal(err)
	} else if len(a) != 2 {
		t.Fatalf("unexpected event count(bob): %d", len(a))
	} else if a[0] != (edb.Event{ID: "1", Type: "PushEvent", Timestamp: MustParseTime("2000-01-01T00:00:00Z"), Actor: "bob"}) {
		t.Fatalf("unexpected event(0): %#v", a[0])
	} else if a[1] != (edb.Event{ID: "2", Type: "IssueEvent", Timestamp: MustParseTime("2000-02-01T00:00:00Z"), Actor: "bob"}) {
		t.Fatalf("unexpected event(1): %#v", a[1])
	}
}

// Ensure all events retrieved from the database.
func TestDB_Events(t *testing.T) {
	db := edb.NewDB()
	if err := db.Open(MustTempFile()); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Save events.
	if err := db.SaveEvents([]edb.Event{
		{ID: "1", Type: "PushEvent", Timestamp: MustParseTime("2000-01-01T00:00:00Z"), Actor: "bob"},
		{ID: "2", Type: "IssueEvent", Timestamp: MustParseTime("2000-02-01T00:00:00Z"), Actor: "bob"},
		{ID: "3", Type: "PullEvent", Timestamp: MustParseTime("2000-01-06T00:00:00Z"), Actor: "susy"},
	}); err != nil {
		t.Fatalf("save events: %s", err)
	}

	// Retrieve events.
	a, err := db.Events()
	if err != nil {
		t.Fatal(err)
	} else if len(a) != 3 {
		t.Fatalf("unexpected event count(bob): %d", len(a))
	} else if a[0] != (edb.Event{ID: "1", Type: "PushEvent", Timestamp: MustParseTime("2000-01-01T00:00:00Z"), Actor: "bob"}) {
		t.Fatalf("unexpected event(0): %#v", a[0])
	} else if a[1] != (edb.Event{ID: "3", Type: "PullEvent", Timestamp: MustParseTime("2000-01-06T00:00:00Z"), Actor: "susy"}) {
		t.Fatalf("unexpected event(1): %#v", a[1])
	} else if a[2] != (edb.Event{ID: "2", Type: "IssueEvent", Timestamp: MustParseTime("2000-02-01T00:00:00Z"), Actor: "bob"}) {
		t.Fatalf("unexpected event(2): %#v", a[2])
	}
}

// MustTempFile returns a temporary path.
func MustTempFile() string {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err.Error())
	}
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func MustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err.Error())
	}
	return t
}
