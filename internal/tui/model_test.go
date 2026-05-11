package tui

import (
	"testing"

	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

func TestMatchesSearchesImportantFields(t *testing.T) {
	entry := scanner.PortEntry{
		Protocol:  "tcp",
		Address:   "127.0.0.1:8080",
		Process:   "node",
		Command:   "node server.js",
		Container: "web",
		Image:     "example/web",
	}

	for _, filter := range []string{"8080", "node", "server", "web", "tcp"} {
		if !matches(entry, filter) {
			t.Fatalf("expected %q to match", filter)
		}
	}
	if matches(entry, "postgres") {
		t.Fatal("unexpected match")
	}
}
