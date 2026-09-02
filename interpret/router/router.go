// Package router selects the most appropriate GoKit interpreter for human input.
package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/date"
	"github.com/codedbyhassan/gokit/interpret/expression"
	"github.com/codedbyhassan/gokit/interpret/number"
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
// The router is deterministic: it collects valid interpretations and selects
// the highest-confidence result. Equal-confidence interpretations are rejected
// rather than silently guessing.
func Parse(input string) (Result, error) {
	original := input
	if strings.TrimSpace(input) == "" {
		return Result{OriginalInput: original}, interpret.ErrEmptyInput
	}

	candidates := make([]Result, 0, 3)

	if result, err := expression.Parse(input); err == nil {
		candidates = append(candidates, Result{
			Kind:          Expression,
			Value:         result.Value,
			Format:        result.Format,
			Confidence:    result.Confidence,
			Ambiguous:     result.Ambiguous,
			Assumptions:   result.Assumptions,
			Candidates:    result.Candidates,
			OriginalInput: original,
		})
	}

	if result, err := date.Parse(input); err == nil {
		candidates = append(candidates, Result{
			Kind:          Date,
			Value:         result.Value,
			Format:        result.Format,
			Confidence:    result.Confidence,
			Ambiguous:     result.Ambiguous,
			Assumptions:   result.Assumptions,
			Candidates:    result.Candidates,
			OriginalInput: original,
		})
	}

	if result, err := number.Parse(input); err == nil {
		candidates = append(candidates, Result{
			Kind:          Number,
			Value:         result.Value,
			Format:        result.Format,
			Confidence:    result.Confidence,
			Ambiguous:     result.Ambiguous,
			Assumptions:   result.Assumptions,
			Candidates:    result.Candidates,
			OriginalInput: original,
		})
	}

	if len(candidates) == 0 {
		return Result{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
	}

	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.Confidence > best.Confidence {
			best = candidate
		}
	}

	var tied []Result
	for _, candidate := range candidates {
		if candidate.Confidence == best.Confidence {
			tied = append(tied, candidate)
		}
	}
	if len(tied) > 1 {
		result := Result{Kind: "ambiguous", Format: "ambiguous", Confidence: best.Confidence, Ambiguous: true, OriginalInput: original}
		for _, candidate := range tied {
			result.Candidates = append(result.Candidates, interpret.Candidate{
				Format:     string(candidate.Kind),
				Confidence: candidate.Confidence,
				Reason:     candidate.Format,
			})
		}
		return result, interpret.ErrAmbiguousInput
	}

	return best, nil
}

// AsNumber returns the parsed number when the result is a number.
func (r Result) AsNumber() (float64, bool) {
	if r.Kind != Number {
		return 0, false
	}
	value, ok := r.Value.(float64)
	return value, ok
}

// AsDate returns the parsed date when the result is a date.
func (r Result) AsDate() (time.Time, bool) {
	if r.Kind != Date {
		return time.Time{}, false
	}
	value, ok := r.Value.(time.Time)
	return value, ok
}

// AsExpression returns the parsed expression when the result is an expression.
func (r Result) AsExpression() (expression.Operation, bool) {
	if r.Kind != Expression {
		return expression.Operation{}, false
	}
	value, ok := r.Value.(expression.Operation)
	return value, ok
}
