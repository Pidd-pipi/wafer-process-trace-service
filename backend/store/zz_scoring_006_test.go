package store

import "testing"

func TestLotListNoPanic(t *testing.T) {
	s := New()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("List panicked: %v", r)
		}
	}()
	items := s.List()
	if len(items) != 2 {
		t.Fatalf("expected 2 seed lots, got %d", len(items))
	}
}

func TestLotUpdateMissingKeepsNotFound(t *testing.T) {
	s := New()
	_, err := s.UpdateStatus("MISSING", "hold", "now")
	if err == nil {
		t.Fatal("UpdateStatus swallowed the not-found error")
	}
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
