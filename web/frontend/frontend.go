// Package frontend provides GoKit's server-rendered web interface.
package frontend

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/codedbyhassan/gokit/interpret/pipeline"
)

//go:embed templates/*.html static/*
var assets embed.FS

type PageData struct {
	Input    string
	Error    string
	Result   *pipeline.Result
	Examples []string
}

type Server struct {
	template *template.Template
}

func NewServer() (*Server, error) {
	tmpl, err := template.New("index.html").Funcs(template.FuncMap{
		"mul": func(a float64, b float64) float64 { return a * b },
		"add": func(a int, b int) int { return a + b },
	}).ParseFS(assets, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{template: tmpl}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("POST /", s.interpret)
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, PageData{Examples: examples()})
}

func (s *Server) interpret(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.render(w, http.StatusBadRequest, PageData{Error: "invalid form submission", Examples: examples()})
		return
	}
	input := strings.TrimSpace(r.FormValue("input"))
	data := PageData{Input: input, Examples: examples()}
	result, err := pipeline.Parse(input)
	if err != nil {
		data.Error = err.Error()
		s.render(w, http.StatusUnprocessableEntity, data)
		return
	}
	data.Result = &result
	s.render(w, http.StatusOK, data)
}

func (s *Server) render(w http.ResponseWriter, status int, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = s.template.ExecuteTemplate(w, "index.html", data)
}

func examples() []string {
	return []string{"what is 20% of 500", "1.5k", "25 km to miles", "how old is someone born 11-11-2011"}
}
