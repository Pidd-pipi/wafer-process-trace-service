package main

import "testing"

func TestProcessBatchReleasesHandles(t *testing.T) {
	m := newOpsHandleManager(16)
	var openDuring int
	err := m.ProcessBatch([]string{"a", "b", "c"}, func(id string) error {
		if id == "c" {
			openDuring = m.OpenCount()
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ProcessBatch returned error: %v", err)
	}
	if openDuring != 1 {
		t.Fatalf("expected only current handle open, got %d (defer accumulated in loop)", openDuring)
	}
}

func TestReleaseFreesSlot(t *testing.T) {
	m := newOpsHandleManager(16)
	if _, err := m.Acquire("x"); err != nil {
		t.Fatal(err)
	}
	if err := m.Release("x"); err != nil {
		t.Fatal(err)
	}
	if got := m.OpenCount(); got != 0 {
		t.Fatalf("release did not free slot: open count %d", got)
	}
}

func TestAcquireRejectsDuplicate(t *testing.T) {
	m := newOpsHandleManager(16)
	if _, err := m.Acquire("x"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Acquire("x"); err == nil {
		t.Fatal("duplicate acquire should be rejected")
	}
}
