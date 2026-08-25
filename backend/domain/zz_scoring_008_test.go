package domain

import "testing"

func TestLotIsPaused(t *testing.T) {
	l := Lot{Status: "paused"}
	if !l.IsPaused() {
		t.Fatal("paused lot reported as not paused")
	}
}

func TestLotIsActiveIncludesPaused(t *testing.T) {
	l := Lot{Status: "paused"}
	if !l.IsActive() {
		t.Fatal("paused lot not considered active")
	}
	r := Lot{Status: "running"}
	if !r.IsActive() {
		t.Fatal("running lot not considered active")
	}
}
