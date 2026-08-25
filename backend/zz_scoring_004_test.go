package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestStateMachineConcurrentMoveAndRead(t *testing.T) {
	m := newOpsStateMachine()
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			_ = m.Move(OpsStatusQueued, OpsStatusActive, "operator")
			_ = m.Move(OpsStatusActive, OpsStatusClosed, "operator")
		}()
		go func() {
			defer wg.Done()
			<-start
			_, _ = m.Last()
			_ = m.History()
			m.Reset()
		}()
	}
	close(start)
	wg.Wait()
}

func TestMiddlewareConcurrentLatency(t *testing.T) {
	h := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/lots", nil))
		}()
	}
	close(start)
	wg.Wait()
}
