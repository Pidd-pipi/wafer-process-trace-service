package health

import (
	"encoding/json"
	"net/http"
)

func Handler(service string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		payload := map[string]string{"status": "ok", "service": service}
		_ = json.NewEncoder(w).Encode(payload)
	}
}
