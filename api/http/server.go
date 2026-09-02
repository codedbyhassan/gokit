// Package http provides a small HTTP/JSON adapter for GoKit's pipeline.
package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/pipeline"
)

type Server struct{}

func NewServer() *Server { return &Server{} }

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /v1/interpret", s.interpret)
	return mux
}

type interpretRequest struct { Input string `json:"input"` }
type interpretResponse struct {
	OriginalInput string `json:"original_input"`
	Source string `json:"source"`
	Kind string `json:"kind"`
	Value any `json:"value"`
	Confidence interpret.Confidence `json:"confidence"`
	Assumptions []string `json:"assumptions,omitempty"`
	Plan pipeline.Result `json:"-"`
}
type errorResponse struct { Error string `json:"error"` }

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status":"ok"})
}

func (s *Server) interpret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { writeError(w, http.StatusMethodNotAllowed, "method not allowed"); return }
	defer r.Body.Close()
	var req interpretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, http.StatusBadRequest, "invalid JSON request"); return }
	result, err := pipeline.Parse(req.Input)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if errors.Is(err, interpret.ErrEmptyInput) { status = http.StatusBadRequest }
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, struct {
		OriginalInput string `json:"original_input"`
		Source pipeline.Source `json:"source"`
		Kind string `json:"kind"`
		Value any `json:"value"`
		Confidence interpret.Confidence `json:"confidence"`
		Assumptions []string `json:"assumptions,omitempty"`
		Plan any `json:"plan"`
	}{result.OriginalInput,result.Source,result.Kind,result.Value,result.Confidence,result.Assumptions,result.Plan})
}

func writeJSON(w http.ResponseWriter, status int, value any) { w.Header().Set("Content-Type","application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(value) }
func writeError(w http.ResponseWriter, status int, message string) { writeJSON(w,status,errorResponse{Error:message}) }
