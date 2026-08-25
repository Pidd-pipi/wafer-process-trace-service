package main

import (
	"errors"
	"sync"
)

var (
	errOpsHandleLimit   = errors.New("operations handle limit exceeded")
	errOpsHandleBusy    = errors.New("operations handle already open")
	errOpsHandleMissing = errors.New("operations handle not open")
)

// OpsHandle simulates an exclusive operating handle for one record. Every
// acquired handle must eventually be released or the slot stays occupied.
type OpsHandle struct {
	id   string
	done bool
}

// OpsHandleManager enforces a fixed limit of concurrently open handles.
type OpsHandleManager struct {
	mu      sync.Mutex
	open    map[string]*OpsHandle
	opened  int
	released int
	maxOpen int
}

func newOpsHandleManager(maxOpen int) *OpsHandleManager {
	if maxOpen < 1 {
		maxOpen = 16
	}
	return &OpsHandleManager{open: map[string]*OpsHandle{}, maxOpen: maxOpen}
}

// Acquire reserves a handle slot for id.
func (m *OpsHandleManager) Acquire(id string) (*OpsHandle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.open) >= m.maxOpen {
		return nil, errOpsHandleLimit
	}
	if _, exists := m.open[id]; exists {
		return nil, errOpsHandleBusy
	}
	h := &OpsHandle{id: id}
	m.open[id] = h
	m.opened++
	return h, nil
}

// Release returns a previously acquired handle slot.
func (m *OpsHandleManager) Release(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	h, ok := m.open[id]
	if !ok {
		return errOpsHandleMissing
	}
	h.done = true
	delete(m.open, id)
	m.released++
	return nil
}

// OpenCount reports how many handle slots are currently occupied.
func (m *OpsHandleManager) OpenCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.open)
}

// Totals reports cumulative acquire/release counts.
func (m *OpsHandleManager) Totals() (opened, released int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.opened, m.released
}

// ProcessBatch runs fn for every id, acquiring and releasing one handle each.
func (m *OpsHandleManager) ProcessBatch(ids []string, fn func(id string) error) error {
	for _, id := range ids {
		h, err := m.Acquire(id)
		if err != nil {
			return err
		}
		if err := fn(id); err != nil {
			_ = m.Release(id)
			return err
		}
		_ = h
		if err := m.Release(id); err != nil {
			return err
		}
	}
	return nil
}
