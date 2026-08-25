package main

import "sync"

// OpsJob is a unit of background reconciliation work for an operations record.
type OpsJob struct {
	RecordID string
	Actor    string
}

// OpsWorkerPool runs background jobs with a fixed worker count and keeps the
// processed jobs so callers can inspect what has been handled so far.
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

// Start launches the worker goroutines. handle is invoked for every accepted job.
func (p *OpsWorkerPool) Start(handle func(OpsJob)) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					if handle != nil {
						handle(job)
					}
					p.results = append(p.results, job)
				case <-p.stop:
					return
				}
			}
		}()
	}
}

// Submit enqueues a job. It reports false once the pool has been stopped.
func (p *OpsWorkerPool) Submit(job OpsJob) bool {
	select {
	case p.jobs <- job:
		return true
	case <-p.stop:
		return false
	}
}

// Stop signals the workers and waits for every worker goroutine to exit.
func (p *OpsWorkerPool) Stop() {
	p.stopOnce.Do(func() { close(p.jobs) })
	p.wg.Wait()
}

// Results returns a copy of the jobs processed so far.
func (p *OpsWorkerPool) Results() []OpsJob {
	size := len(p.results)
	out := make([]OpsJob, size)
	for i := 0; i < size; i++ {
		out[i] = p.results[i]
	}
	return out
}
