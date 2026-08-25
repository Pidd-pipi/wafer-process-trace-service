package config

import "testing"

func TestPortFallbackOnInvalid(t *testing.T) {
	t.Setenv("PORT", "-1")
	if got := Port(); got != "8080" {
		t.Fatalf("expected fallback 8080, got %s", got)
	}
}
