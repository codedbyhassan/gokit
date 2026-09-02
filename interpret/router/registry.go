package router

import (
	"sort"

	"github.com/codedbyhassan/gokit/interpret"
)

// Interpreter is a domain-specific parser that can participate in routing.
// Implementations should return an error when the input does not belong to
// their domain.
type Interpreter interface {
	Interpret(input string) (Candidate, error)
}

// Candidate is a router-ready interpretation produced by an Interpreter.
type Candidate struct {
	Kind        Kind
	Value       any
	Format      string
	Confidence  interpret.Confidence
	Ambiguous   bool
	Assumptions []string
	Candidates  []interpret.Candidate
}

// Registry stores interpreters in deterministic priority order.
type Registry struct {
	interpreters []Interpreter
}

// NewRegistry creates an empty interpreter registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds an interpreter to the registry. Registration order is retained
// only as a deterministic tie-breaker; confidence remains the primary signal.
func (r *Registry) Register(interpreter Interpreter) {
	if interpreter != nil {
		r.interpreters = append(r.interpreters, interpreter)
	}
}

func (r *Registry) candidates(input string) []Candidate {
	out := make([]Candidate, 0, len(r.interpreters))
	for _, interpreter := range r.interpreters {
		candidate, err := interpreter.Interpret(input)
		if err == nil {
			out = append(out, candidate)
		}
	}
	return out
}

func rankCandidates(candidates []Candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Confidence > candidates[j].Confidence
	})
}
