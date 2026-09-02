package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	NewServer().Handler().ServeHTTP(r, req)
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), `"status":"ok"`) {
		t.Fatalf("status=%d body=%s", r.Code, r.Body.String())
	}
}

func TestInterpret(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/interpret", strings.NewReader(`{"input":"what is 20% of 500"}`))
	req.Header.Set("Content-Type", "application/json")
	NewServer().Handler().ServeHTTP(r, req)

	if r.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", r.Code, r.Body.String())
	}

	var body struct {
		Kind  string `json:"kind"`
		Value struct {
			Left      float64 `json:"Left"`
			Right     float64 `json:"Right"`
			Operation string  `json:"Operation"`
			Value     float64 `json:"Value"`
		} `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Kind != "calculate" {
		t.Fatalf("expected calculate kind, got %q", body.Kind)
	}
	if body.Value.Value != 100 {
		t.Fatalf("expected value=100, got %v", body.Value.Value)
	}
}

func TestInterpretRejectsInvalidJSON(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/interpret", strings.NewReader("{"))
	NewServer().Handler().ServeHTTP(r, req)
	if r.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", r.Code)
	}
}
