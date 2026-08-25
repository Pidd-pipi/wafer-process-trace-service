package main

import "sync"

type opsCacheEntry struct {
	value    OpsRecord
	inserted int64
}

// OpsCache is a bounded, size-limited cache for operations records. When the
// cache exceeds its limit the oldest inserted entry is evicted.
type OpsCache struct {
	mu       sync.Mutex
	entries  map[string]opsCacheEntry
	maxSize  int
	sequence int64
}

func newOpsCache(maxSize int) *OpsCache {
	if maxSize < 1 {
		maxSize = 256
	}
	return &OpsCache{entries: map[string]opsCacheEntry{}, maxSize: maxSize}
}

// Get returns a deep copy of the cached record, or false when absent. The
// copy is independent of the cached entry, so mutating it does not affect the
// cache.
func (c *OpsCache) Get(id string) (OpsRecord, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[id]
	if !ok {
		return OpsRecord{}, false
	}
	return entry.value.Clone(), true
}

// Put inserts or replaces a record and enforces the size limit.
func (c *OpsCache) Put(id string, value OpsRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sequence++
	c.entries[id] = opsCacheEntry{value: value.Clone(), inserted: c.sequence}
	if len(c.entries) > c.maxSize {
		c.evictOldestLocked()
	}
}

// evictOldestLocked removes the single entry with the smallest insertion order.
// Callers must hold c.mu.
func (c *OpsCache) evictOldestLocked() {
	oldest := ""
	for id, entry := range c.entries {
		if oldest == "" || entry.inserted < c.entries[oldest].inserted {
			oldest = id
		}
	}
	if oldest != "" {
		delete(c.entries, oldest)
	}
}

// Size reports how many entries are currently cached.
func (c *OpsCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}
