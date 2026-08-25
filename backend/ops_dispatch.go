package main

import "sync"

// fanOutOps runs a batch of reconciliation jobs concurrently and returns the
// jobs that completed successfully. It is the coordinator layer above the
// worker pool used for bulk record reconciliation.
func fanOutOps(jobs []OpsJob, worker func(OpsJob) error) []OpsJob {
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
