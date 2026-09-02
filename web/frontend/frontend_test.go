package frontend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerGET(t *testing.T) {
	server, err := NewServer()
	if err != nil { t.Fatal(err) }
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status=%d want %d", rec.Code, http.StatusOK) }
	if !strings.Contains(rec.Body.String(), "What can GoKit") { t.Fatal("expected workspace heading") }
	if !strings.Contains(rec.Body.String(), "gokit-local") { t.Fatal("expected local-first frontend") }
}

func TestServerPOST(t *testing.T) {
	server, err := NewServer()
	if err != nil { t.Fatal(err) }
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("input=what+is+20%25+of+500"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status=%d want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String()) }
	if !strings.Contains(rec.Body.String(), "100") { t.Fatal("expected calculated value in response") }
}

func TestServerAPI(t *testing.T) {
	server, err := NewServer()
	if err != nil { t.Fatal(err) }
	req := httptest.NewRequest(http.MethodPost, "/v1/interpret", strings.NewReader(`{"input":"what is 20% of 500"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status=%d want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String()) }
	var body struct { Kind string `json:"kind"`; Value any `json:"value"` }
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil { t.Fatal(err) }
	if body.Kind != "calculate" { t.Fatalf("kind=%q want calculate", body.Kind) }
	if body.Value != float64(100) { t.Fatalf("value=%v want 100", body.Value) }
}

func TestServerHealth(t *testing.T) {
	server, err := NewServer()
	if err != nil { t.Fatal(err) }
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status=%d want %d", rec.Code, http.StatusOK) }
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) { t.Fatal("expected health response") }
}
