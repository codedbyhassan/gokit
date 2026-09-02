// Package router selects the most appropriate GoKit interpreter for human input.
package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/expression"
)

// Kind identifies the domain selected by the router.
type Kind string

const (
	Number     Kind = "number"
	Date       Kind = "date"
	Expression Kind = "expression"
)

// Result is the unified output of the interpreter router.
type Result struct {
	Kind          Kind
	Value         any
	Format        string
	Confidence    interpret.Confidence
	Ambiguous     bool
	Assumptions   []string
	Candidates    []interpret.Candidate
	OriginalInput string
}

// Parse determines which built-in interpreter best matches the input.
func Parse(input string) (Result, error) {
	return DefaultRegistry().Parse(input)
}

// Parse runs the registered interpreters and selects the highest-confidence
// interpretation. Equal-confidence matches are rejected rather than guessed.
func (r *Registry) Parse(input string) (Result, error) {
	original := input
	if strings.TrimSpace(input) == "" {
		return Result{OriginalInput: original}, interpret.ErrEmptyInput
	}

	candidates := r.candidates(input)
	if len(candidates) == 0 {
		return Result{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
	}

	rankCandidates(candidates)
	best := candidates[0]
	if len(candidates) > 1 && candidates[1].Confidence == best.Confidence {
		result := Result{Kind: "ambiguous", Format: "ambiguous", Confidence: best.Confidence, Ambiguous: true, OriginalInput: original}
		for _, candidate := range candidates {
			if candidate.Confidence != best.Confidence {
				break
			}
			result.Candidates = append(result.Candidates, interpret.Candidate{
				Format:     string(candidate.Kind),
				Confidence: candidate.Confidence,
				Reason:     candidate.Format,
			})
		}
		return result, interpret.ErrAmbiguousInput
	}

	return Result{
		Kind:          best.Kind,
		Value:         best.Value,
		Format:        best.Format,
		Confidence:    best.Confidence,
		Ambiguous:     best.Ambiguous,
		Assumptions:   best.Assumptions,
		Candidates:    best.Candidates,
		OriginalInput: original,
	}, nil
}

// AsNumber returns the parsed number when the result is a number.
func (r Result) AsNumber() (float64, bool) {
	if r.Kind != Number { return 0, false }
	value, ok := r.Value.(float64)
	return value, ok
}

// AsDate returns the parsed date when the result is a date.
func (r Result) AsDate() (time.Time, bool) {
	if r.Kind != Date { return time.Time{}, false }
	value, ok := r.Value.(time.Time)
	return value, ok
}

// AsExpression returns the parsed expression when the result is an expression.
func (r Result) AsExpression() (expression.Operation, bool) {
	if r.Kind != Expression { return expression.Operation{}, false }
	value, ok := r.Value.(expression.Operation)
	return value, ok
}
