package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestWorkerPoolConcurrentResults(t *testing.T) {
	p := newOpsWorkerPool(8)
	gate := make(chan struct{})
	p.Start(func(j OpsJob) { <-gate })
	for i := 0; i < 8; i++ {
		p.Submit(OpsJob{RecordID: fmt.Sprintf("r%d", i)})
	}
	close(gate)
	p.Stop()
	if got := len(p.Results()); got != 8 {
		t.Fatalf("expected 8 processed jobs, got %d", got)
	}
}

func TestSubmitAfterStopNoPanic(t *testing.T) {
	p := newOpsWorkerPool(2)
	p.Start(func(j OpsJob) {})
	p.Submit(OpsJob{RecordID: "r1"})
	p.Stop()
	accepted := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Submit panicked after Stop: %v", r)
			}
		}()
		accepted = p.Submit(OpsJob{RecordID: "r2"})
	}()
	if accepted {
		t.Fatal("Submit accepted a job after Stop")
	}
}

func TestFanOutConcurrentCounter(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			fanOutOps([]OpsJob{{RecordID: "a"}, {RecordID: "b"}}, func(j OpsJob) error { return nil })
		}()
	}
	close(start)
	wg.Wait()
}
