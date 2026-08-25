package main

import (
	"sync"
	"sync/atomic"
)

var opsFanOutJobs int64

// fanOutOps runs a batch of reconciliation jobs concurrently and returns the
// jobs that completed successfully. It is the coordinator layer above the
// worker pool used for bulk record reconciliation.
//
// Concurrency notes:
//   - opsFanOutJobs is shared across batches, so the update is atomic.
//   - Each job runs on its own goroutine and reports success on a buffered
//     results channel; the close(results) after Wait guarantees the range
//     below terminates. errs is buffered to len(jobs) so worker error sends
//     never block even if the caller never drains them.
//   - Submit/Stop on the underlying pool is the caller's job; fanOutOps only
//     fans the per-job work out.
func fanOutOps(jobs []OpsJob, worker func(OpsJob) error) []OpsJob {
	atomic.AddInt64(&opsFanOutJobs, int64(len(jobs)))
	var wg sync.WaitGroup
	results := make(chan OpsJob, len(jobs))
	errs := make(chan error, len(jobs))
	for _, job := range jobs {
		wg.Add(1)
		go func(j OpsJob) {
			defer wg.Done()
			if err := worker(j); err != nil {
				errs <- err
				return
			}
			results <- j
		}(job)
	}
	wg.Wait()
	close(results)
	close(errs)
	out := make([]OpsJob, 0, len(jobs))
	for j := range results {
		out = append(out, j)
	}
	return out
}
