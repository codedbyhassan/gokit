package frontend

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServerGET(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Talk to your Go utilities") {
		t.Fatal("expected frontend heading")
	}
}

func TestServerPOST(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("input=what+is+20%25+of+500"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "100") {
		t.Fatal("expected calculated value in response")
	}
}
