package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/wafer-process-trace-service/store"
	"example.com/wafer-process-trace-service/web"
)

func TestPausedStatusAccepted(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/lots/status", strings.NewReader(`{"lot_id":"LOT-24082","status":"paused"}`)))
	if post.Code != http.StatusOK {
		t.Fatalf("expected 200 for paused status, got %d: %s", post.Code, post.Body.String())
	}
}

func TestPausedLotsRoute(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/lots/status", strings.NewReader(`{"lot_id":"LOT-24082","status":"completed"}`)))
	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/lots/paused", nil))
	if get.Code != http.StatusOK {
		t.Fatalf("expected 200 from paused route, got %d", get.Code)
	}
	if strings.Contains(get.Body.String(), "LOT-24082") {
		t.Fatalf("completed lot leaked into paused list: %s", get.Body.String())
	}
}
