package main

import (
	"context"
	"errors"
	"testing"
)

func TestOpsMissingRecordErrorCode(t *testing.T) {
	s := newOpsStore(nil)
	_, err := s.Get(context.Background(), "missing-lot")
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("expected ErrOpsNotFound chain, got %v", err)
	}
	if code := opsCode(err); code != "not_found" {
		t.Fatalf("expected not_found, got %s", code)
	}
}

func TestOpsUpdateMissingErrorCode(t *testing.T) {
	s := newOpsStore(nil)
	err := s.Update(context.Background(), OpsRecord{ID: "missing-lot"}, 0)
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("expected ErrOpsNotFound chain, got %v", err)
	}
	if code := opsCode(err); code != "not_found" {
		t.Fatalf("expected not_found, got %s", code)
	}
}

func TestOpsDeleteMissingErrorCode(t *testing.T) {
	s := newOpsStore(nil)
	err := s.Delete(context.Background(), "missing-lot")
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("expected ErrOpsNotFound chain, got %v", err)
	}
	if code := opsCode(err); code != "not_found" {
		t.Fatalf("expected not_found, got %s", code)
	}
}

func TestOpsPutConflictPreservesCode(t *testing.T) {
	s := newOpsStore([]OpsRecord{{ID: "dup-lot", Status: OpsStatusActive}})
	err := s.Put(context.Background(), OpsRecord{ID: "dup-lot"})
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("expected ErrOpsConflict chain, got %v", err)
	}
	if code := opsCode(err); code != "conflict" {
		t.Fatalf("expected conflict, got %s", code)
	}
}

func TestOpsWrappedErrorPreservesCause(t *testing.T) {
	err := wrapOps("create", "store.put", ErrOpsNotFound)
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("expected ErrOpsNotFound chain through wrapOps, got %v", err)
	}
}
