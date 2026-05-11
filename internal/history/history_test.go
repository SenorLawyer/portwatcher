package history

import (
	"testing"
	"time"

	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

func TestDiffOpenedClosedChanged(t *testing.T) {
	at := time.Now()
	prev := scanner.Snapshot{
		At: at,
		Ports: []scanner.PortEntry{
			{ID: "tcp/127.0.0.1/3000/10", Port: 3000, Command: "old", State: "listen"},
			{ID: "tcp/127.0.0.1/4000/11", Port: 4000, Command: "gone", State: "listen"},
		},
	}
	next := scanner.Snapshot{
		At: at.Add(time.Second),
		Ports: []scanner.PortEntry{
			{ID: "tcp/127.0.0.1/3000/10", Port: 3000, Command: "new", State: "listen"},
			{ID: "tcp/127.0.0.1/5000/12", Port: 5000, Command: "added", State: "listen"},
		},
	}

	events := Diff(prev, next)
	counts := map[EventType]int{}
	for _, event := range events {
		counts[event.Type]++
	}

	if counts[Changed] != 1 || counts[Closed] != 1 || counts[Opened] != 1 {
		t.Fatalf("unexpected events: %#v", events)
	}
}
