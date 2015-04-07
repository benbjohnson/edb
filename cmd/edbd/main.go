package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"code.google.com/p/goauth2/oauth"
	"github.com/BurntSushi/toml"
	"github.com/benbjohnson/edb"
	"github.com/google/go-github/github"
)

var (
	ErrConfigRequired = errors.New("config required")

	ErrDataPathRequired = errors.New("data path required")

	ErrAccessTokenRequired = errors.New("access token required")
)

func main() {
	m := NewMain()
	defer func() { _ = m.Close() }()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-(chan struct{})(nil)
}

// Main represents the main program execution.
type Main struct {
	logger *log.Logger

	closing chan struct{}
	wg      sync.WaitGroup

	DB *edb.DB

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		closing: make(chan struct{}),
	}
}

func (m *Main) Close() error {
	if m.DB != nil {
		_ = m.DB.Close()
	}
	if m.closing != nil {
		close(m.closing)
		m.closing = nil
	}
	m.wg.Wait()
	return nil
}

func (m *Main) Run(args ...string) error {
	// Set up logger.
	m.logger = log.New(m.Stderr, "", log.LstdFlags)

	// Parse command line flags and config.
	config, err := m.parseFlags(args)
	if err != nil {
		return err
	}

	m.logger.Println("starting up...")

	// Open database.
	m.DB = edb.NewDB()
	if err := m.DB.Open(config.DataPath); err != nil {
		return err
	}

	// Start fetchers.
	m.startFetchers(config.AccessToken, config.Usernames)

	return nil
}

func (m *Main) parseFlags(args []string) (*Config, error) {
	// Parse command line flags.
	fs := flag.NewFlagSet("edb", flag.ContinueOnError)
	configPath := fs.String("config", "", "config path")
	if err := fs.Parse(args); err != nil {
		return nil, err
	} else if *configPath == "" {
		return nil, ErrConfigRequired
	}

	// Parse configuration.
	var config Config
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		return nil, err
	} else if config.DataPath == "" {
		return nil, ErrDataPathRequired
	} else if config.AccessToken == "" {
		return nil, ErrAccessTokenRequired
	}

	return &config, nil
}

// startFetchers starts a fetcher for each of the users.
func (m *Main) startFetchers(accessToken string, usernames []string) {
	// Create GitHub client.
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: accessToken}}
	client := github.NewClient(t.Client())

	m.logger.Printf("starting fetchers(%d)", len(usernames))

	// Start a fetcher for each username
	for _, username := range usernames {
		f := edb.NewGitHubFetcher(client, m.DB, username)
		f.Logger = m.logger

		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			f.Run(m.closing)
		}()
	}
}

// Config represents the application configuration format.
type Config struct {
	AccessToken string   `toml:"access-token"`
	DataPath    string   `toml:"data-path"`
	Usernames   []string `toml:"usernames"`
}
