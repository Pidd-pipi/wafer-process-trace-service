package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"example.com/wafer-process-trace-service/domain"
	"example.com/wafer-process-trace-service/store"
	"example.com/wafer-process-trace-service/validation"
)

type server struct{ store *store.Store }

func (s *server) collection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string][]domain.Lot{"items": s.store.List()})
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	defer r.Body.Close()
	var request struct {
		ID     string `json:"lot_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&request); err != nil || strings.TrimSpace(request.ID) == "" {
		writeError(w, http.StatusBadRequest, "lot_id and status are required")
		return
	}
	if err := validation.Status(request.Status); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := s.store.UpdateStatus(request.ID, request.Status, time.Now().UTC().Format(time.RFC3339))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "lot not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
