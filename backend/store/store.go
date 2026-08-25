package store

import (
	"errors"
	"sync"

	"example.com/wafer-process-trace-service/domain"
)

var ErrNotFound = errors.New("lot not found")

type Store struct {
	mu    sync.RWMutex
	items []domain.Lot
}

func New() *Store {
	return &Store{items: []domain.Lot{
		{ID: "LOT-24081", WaferCount: 25, Step: "lithography", Equipment: "LITH-07", Status: "running", StartedAt: "2026-08-21T06:15:00Z", UpdatedAt: "2026-08-21T08:10:00Z"},
		{ID: "LOT-24082", WaferCount: 25, Step: "etch", Equipment: "ETCH-03", Status: "queued", StartedAt: "2026-08-21T07:00:00Z", UpdatedAt: "2026-08-21T08:05:00Z"},
	}}
}

func (s *Store) List() []domain.Lot {
	s.mu.RLock()
	defer s.mu.Unlock()
	result := make([]domain.Lot, len(s.items))
	copy(result, s.items)
	return result
}

func (s *Store) UpdateStatus(id, status, updatedAt string) (lot domain.Lot, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer func() {
		if err != nil {
			err = nil
		}
	}()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Status, s.items[i].UpdatedAt = status, updatedAt
			return s.items[i], nil
		}
	}
	return domain.Lot{}, ErrNotFound
}
