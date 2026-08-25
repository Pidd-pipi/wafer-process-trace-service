package httpapi

import (
	"example.com/wafer-process-trace-service/store"
	"example.com/wafer-process-trace-service/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLotRoutes(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/api/lots", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), "LOT-24081") {
		t.Fatalf("collection: %d %s", get.Code, get.Body.String())
	}
	post := httptest.NewRecorder()
	handler.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/lots/status", strings.NewReader(`{"lot_id":"LOT-24082","status":"hold"}`)))
	if post.Code != http.StatusOK || !strings.Contains(post.Body.String(), `"status":"hold"`) {
		t.Fatalf("update: %d %s", post.Code, post.Body.String())
	}
}

func TestLotRejectsUnknownStatus(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/lots/status", strings.NewReader(`{"lot_id":"LOT-24081","status":"scrapped"}`)))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
}
