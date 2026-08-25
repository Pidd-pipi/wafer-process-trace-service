package health

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthStatusOk(t *testing.T) {
	rec := httptest.NewRecorder()
	Handler("wafer-process-trace-service")(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("health status not ok: %s", rec.Body.String())
	}
}
