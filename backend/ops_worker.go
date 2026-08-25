package main

import (
	"sync"
)

// OpsJob is a unit of background reconciliation work for an operations record.
type OpsJob struct {
	RecordID string
	Actor    string
}

// OpsWorkerPool runs background jobs with a fixed worker count and keeps the
// processed jobs so callers can inspect what has been handled so far.
//
// Lifecycle is strictly Start -> (Submit|Results|Stop)* -> Stop. The jobs
// channel is intentionally NEVER closed: closing it would let Stop race with
// an in-flight Submit, and a Go select can non-deterministically pick the
// send-on-closed case, panicking. Instead, workers exit when the stop channel
// is closed, and Submit aborts on the same signal — so Submit after Stop
// returns false instead of panicking. The processed-results slice is guarded
// by the pool mutex so concurrent worker appends and Results reads are
// race-free.
type OpsWorkerPool struct {
	mu       sync.Mutex
	workers  int
	jobs     chan OpsJob
	results  []OpsJob
	stop     chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

func newOpsWorkerPool(workers int) *OpsWorkerPool {
	if workers < 1 {
		workers = 1
	}
	return &OpsWorkerPool{
		workers: workers,
		jobs:    make(chan OpsJob),
		stop:    make(chan struct{}),
	}
}

// Start launches the worker goroutines. handle is invoked for every accepted
// job. A worker never exits mid-handle: the stop signal is only observed at
// the top of the select loop, so any job that was accepted is fully processed
// and recorded before the worker tears down — no early termination.
func (p *OpsWorkerPool) Start(handle func(OpsJob)) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case job := <-p.jobs:
					if handle != nil {
						handle(job)
					}
					p.mu.Lock()
					p.results = append(p.results, job)
					p.mu.Unlock()
				case <-p.stop:
					return
				}
			}
		}()
	}
}

// Submit enqueues a job. It blocks until either the job is accepted by a worker
// or the pool has been stopped. After Stop it reports false and never blocks or
// panics. The non-blocking stop check up front keeps the common post-shutdown
// path lock-free; the second select handles the concurrent Stop case.
func (p *OpsWorkerPool) Submit(job OpsJob) bool {
	select {
	case <-p.stop:
		return false
	default:
	}
	select {
	case p.jobs <- job:
		return true
	case <-p.stop:
		return false
	}
}

// Stop signals the workers and waits for every worker goroutine to exit.
// stopOnce makes it idempotent, so concurrent or repeated Stop calls are safe.
// The jobs channel is left open on purpose (see the struct comment).
func (p *OpsWorkerPool) Stop() {
	p.stopOnce.Do(func() { close(p.stop) })
	p.wg.Wait()
}

// Results returns a copy of the jobs processed so far.
func (p *OpsWorkerPool) Results() []OpsJob {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]OpsJob, len(p.results))
	copy(out, p.results)
	return out
}
