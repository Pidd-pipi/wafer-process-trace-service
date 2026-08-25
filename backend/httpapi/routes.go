package httpapi

import (
	"io/fs"
	"net/http"

	"example.com/wafer-process-trace-service/health"
	"example.com/wafer-process-trace-service/store"
)

func NewHandler(st *store.Store, staticFS fs.FS) http.Handler {
	s := &server{store: st}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler("wafer-process-trace-service"))
	mux.HandleFunc("/api/lots", s.collection)
	mux.HandleFunc("/api/lots/status", s.status)
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	return mux
}
