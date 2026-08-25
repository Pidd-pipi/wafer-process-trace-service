package main

import (
	"context"
	"testing"
	"time"
)

func TestOpsContextPropagatesCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	child, _ := opsContext(parent, 5*time.Second)
	cancel()
	select {
	case <-child.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("child context did not observe parent cancellation")
	}
}

func TestOpsDeadlineReportsDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !opsDeadline(ctx) {
		t.Fatal("opsDeadline returned false for a context with deadline")
	}
}

func TestOpsDelayCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := opsDelay(ctx, 5*time.Second); err == nil {
		t.Fatal("opsDelay ignored an already cancelled context")
	}
}

func TestShutdownContextHasDeadline(t *testing.T) {
	ctx, cancel := newShutdownContext()
	defer cancel()
	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("shutdown context has no deadline")
	}
}
