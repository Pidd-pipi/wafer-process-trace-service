package httpapi

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"example.com/wafer-process-trace-service/domain"
	"example.com/wafer-process-trace-service/health"
	"example.com/wafer-process-trace-service/store"
)

func NewHandler(st *store.Store, staticFS fs.FS) http.Handler {
	s := &server{store: st}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler("wafer-process-trace-service"))
	mux.HandleFunc("/api/lots", s.collection)
	mux.HandleFunc("/api/lots/status", s.status)
	mux.HandleFunc("/api/lots/paused", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var paused []domain.Lot
		for _, item := range s.store.List() {
			if item.IsPaused() {
				paused = append(paused, item)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string][]domain.Lot{"items": paused})
	})
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	return mux
}
