package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

type EventType string

const (
	Opened  EventType = "opened"
	Closed  EventType = "closed"
	Changed EventType = "changed"
)

type Event struct {
	Type EventType         `json:"type"`
	At   time.Time         `json:"at"`
	Port scanner.PortEntry `json:"port"`
}

type Store struct {
	path      string
	retention time.Duration
	file      *os.File
}

func Open(path string, retention time.Duration) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return &Store{path: path, retention: retention, file: file}, nil
}

func (s *Store) Close() error {
	if s.file == nil {
		return nil
	}
	return s.file.Close()
}

func (s *Store) Append(events []Event) error {
	if len(events) == 0 {
		return nil
	}
	enc := json.NewEncoder(s.file)
	for _, event := range events {
		if err := enc.Encode(event); err != nil {
			return err
		}
	}
	return s.file.Sync()
}

func (s *Store) ReadAll() ([]Event, error) {
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cutoff := time.Now().Add(-s.retention)
	events := []Event{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		var event Event
		if err := json.Unmarshal(sc.Bytes(), &event); err == nil && event.At.After(cutoff) {
			events = append(events, event)
		}
	}
	return events, sc.Err()
}

func Diff(prev, next scanner.Snapshot) []Event {
	if prev.At.IsZero() {
		return nil
	}
	before := index(prev.Ports)
	after := index(next.Ports)
	events := []Event{}

	for id, port := range after {
		old, ok := before[id]
		if !ok {
			events = append(events, Event{Type: Opened, At: next.At, Port: port})
			continue
		}
		if old.Command != port.Command || old.Container != port.Container || old.State != port.State {
			events = append(events, Event{Type: Changed, At: next.At, Port: port})
		}
	}
	for id, port := range before {
		if _, ok := after[id]; !ok {
			events = append(events, Event{Type: Closed, At: next.At, Port: port})
		}
	}
	return events
}

func index(ports []scanner.PortEntry) map[string]scanner.PortEntry {
	out := make(map[string]scanner.PortEntry, len(ports))
	for _, port := range ports {
		out[port.ID] = port
	}
	return out
}
