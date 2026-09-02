package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	NewServer().Handler().ServeHTTP(r, req)
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), `"status":"ok"`) { t.Fatalf("status=%d body=%s", r.Code, r.Body.String()) }
}

func TestInterpret(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/interpret", strings.NewReader(`{"input":"what is 20% of 500"}`))
	req.Header.Set("Content-Type", "application/json")
	NewServer().Handler().ServeHTTP(r, req)
	body := r.Body.String()
	if r.Code != http.StatusOK || !strings.Contains(body, `"kind":"calculate"`) || !strings.Contains(body, `"value":100`) { t.Fatalf("status=%d body=%s", r.Code, body) }
}

func TestInterpretRejectsInvalidJSON(t *testing.T) {
	r := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/interpret", strings.NewReader("{"))
	NewServer().Handler().ServeHTTP(r, req)
	if r.Code != http.StatusBadRequest { t.Fatalf("status=%d", r.Code) }
}
