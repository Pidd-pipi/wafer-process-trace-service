package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/wafer-process-trace-service/store"
	"example.com/wafer-process-trace-service/web"
)

type trackingBody struct {
	*strings.Reader
	closed bool
}

func (b *trackingBody) Close() error {
	b.closed = true
	return nil
}

func TestStatusClosesBody(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	body := &trackingBody{Reader: strings.NewReader(`{"lot_id":"","status":"hold"}`)}
	req := httptest.NewRequest(http.MethodPost, "/api/lots/status", body)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if !body.closed {
		t.Fatal("status handler did not close the request body on the error path")
	}
}

func TestStatusErrorCode(t *testing.T) {
	handler := NewHandler(store.New(), web.FS)
	req := httptest.NewRequest(http.MethodPost, "/api/lots/status", strings.NewReader(`{"lot_id":"LOT-24081","status":"scrapped"}`))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid status, got %d", rec.Code)
	}
	_ = io.Discard
}
