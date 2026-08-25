package main

import (
	"testing"
	"time"
)

func TestAuditForReturnsIsolatedCopy(t *testing.T) {
	a := newOpsAudit()
	a.Add("r1", "created", "alice")
	a.Add("r2", "created", "bob")
	_ = a.For("r2")
	if a.Count() != 2 {
		t.Fatalf("audit event count corrupted: %d", a.Count())
	}
	got := a.For("r1")
	if len(got) != 1 || got[0].RecordID != "r1" {
		t.Fatalf("r1 event corrupted after aliased For: %v", got)
	}
}

func TestAuditSinceReturnsIsolatedCopy(t *testing.T) {
	a := newOpsAudit()
	a.Add("r1", "created", "alice")
	evs := a.Since(time.Now().Add(-time.Hour))
	if len(evs) != 1 {
		t.Fatalf("expected 1 event in window, got %d", len(evs))
	}
	evs[0].Actor = "EVIL"
	got := a.For("r1")
	if len(got) != 1 || got[0].Actor == "EVIL" {
		t.Fatalf("audit mutated through returned slice: %v", got)
	}
}

func TestClonePageIsolatesItems(t *testing.T) {
	p := OpsPage{Items: []OpsRecord{{ID: "a"}, {ID: "b"}}}
	c := opsClonePage(p)
	c.Items[0] = OpsRecord{ID: "mutated"}
	if p.Items[0].ID == "mutated" {
		t.Fatal("clone page shares backing array with source page")
	}
}

func TestSlicePageTrimsCapacity(t *testing.T) {
	items := make([]OpsRecord, 1000)
	for i := range items {
		items[i] = OpsRecord{ID: string(rune('a' + i%26))}
	}
	window := opsSlicePage(items, 10, 20)
	if len(window) != 10 {
		t.Fatalf("expected window length 10, got %d", len(window))
	}
	if cap(window) > len(window) {
		t.Fatalf("window retains oversized backing array: cap %d len %d", cap(window), len(window))
	}
}
