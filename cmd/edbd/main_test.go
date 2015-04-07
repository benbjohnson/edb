package main_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/benbjohnson/edb/cmd/edbd"
)

func TestMain(t *testing.T) {
	m := NewMain()
	defer m.Close()
	if err := m.Run("-config", ""); err != main.ErrConfigRequired {
		t.Fatal(err)
	}
}

func TestConfig_Parse(t *testing.T) {
	s := `
data-path = "/tmp/my.conf"
usernames = [
	"benbjohnson",
	"BurntSushi",
	"worace"
]
`

	var c main.Config
	if _, err := toml.Decode(s, &c); err != nil {
		t.Fatal(err)
	} else if c.DataPath != `/tmp/my.conf` {
		t.Fatalf("unexpected data path: %s", c.DataPath)
	} else if !reflect.DeepEqual(c.Usernames, []string{"benbjohnson", "BurntSushi", "worace"}) {
		t.Fatalf("unexpected usernames: %+v", c.Usernames)
	}
}

// Main represents a test wrapper for main.Main.
type Main struct {
	*main.Main

	Stdin  bytes.Buffer
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func NewMain() *Main {
	m := &Main{
		Main: main.NewMain(),
	}

	m.Main.Stdin = &m.Stdin
	m.Main.Stdout = &m.Stdout
	m.Main.Stderr = &m.Stderr

	return m
}
