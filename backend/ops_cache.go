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

// Get returns a deep copy of the cached record, or false when absent.
func (c *OpsCache) Get(id string) (OpsRecord, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[id]
	if !ok {
		return OpsRecord{}, false
	}
	record := entry.value
	return record, true
}

// Put inserts or replaces a record and enforces the size limit.
func (c *OpsCache) Put(id string, value OpsRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sequence++
	c.entries[id] = opsCacheEntry{value: value.Clone(), inserted: c.sequence}
}

// evictOldestLocked removes the single entry with the smallest insertion order.
// Callers must hold c.mu.
func (c *OpsCache) evictOldestLocked() {
	oldest := ""
	var oldestSeq int64 = -1
	for id, entry := range c.entries {
		if oldestSeq == -1 || entry.inserted > oldestSeq {
			oldest = id
			oldestSeq = entry.inserted
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
