package app

import (
	"context"
	"sync"

	"github.com/SenorLawyer/portwatcher/internal/history"
	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

type Service struct {
	scanner *scanner.Scanner
	store   *history.Store
	mu      sync.Mutex
	prev    scanner.Snapshot
}

func New(scanner *scanner.Scanner, store *history.Store) *Service {
	return &Service{scanner: scanner, store: store}
}

func (s *Service) Refresh(ctx context.Context) (scanner.Snapshot, []history.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	snap, err := s.scanner.Snapshot(ctx)
	if err != nil {
		return scanner.Snapshot{}, nil, err
	}
	events := history.Diff(s.prev, snap)
	if err := s.store.Append(events); err != nil {
		return scanner.Snapshot{}, nil, err
	}
	s.prev = snap
	return snap, events, nil
}

func (s *Service) History() ([]history.Event, error) {
	return s.store.ReadAll()
}
