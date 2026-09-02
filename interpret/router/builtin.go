package router

import (
	"github.com/codedbyhassan/gokit/interpret/date"
	"github.com/codedbyhassan/gokit/interpret/expression"
	"github.com/codedbyhassan/gokit/interpret/number"
	"github.com/codedbyhassan/gokit/interpret/unit"
)

type numberInterpreter struct{}
func (numberInterpreter) Interpret(input string) (Candidate, error) {
	result, err := number.Parse(input)
	if err != nil { return Candidate{}, err }
	return Candidate{Kind: Number, Value: result.Value, Format: result.Format, Confidence: result.Confidence, Ambiguous: result.Ambiguous, Assumptions: result.Assumptions, Candidates: result.Candidates}, nil
}

type dateInterpreter struct{}
func (dateInterpreter) Interpret(input string) (Candidate, error) {
	result, err := date.Parse(input)
	if err != nil { return Candidate{}, err }
	return Candidate{Kind: Date, Value: result.Value, Format: result.Format, Confidence: result.Confidence, Ambiguous: result.Ambiguous, Assumptions: result.Assumptions, Candidates: result.Candidates}, nil
}

type expressionInterpreter struct{}
func (expressionInterpreter) Interpret(input string) (Candidate, error) {
	result, err := expression.Parse(input)
	if err != nil { return Candidate{}, err }
	return Candidate{Kind: Expression, Value: result.Value, Format: result.Format, Confidence: result.Confidence, Ambiguous: result.Ambiguous, Assumptions: result.Assumptions, Candidates: result.Candidates}, nil
}

type unitInterpreter struct{}
func (unitInterpreter) Interpret(input string) (Candidate, error) {
	result, err := unit.Parse(input)
	if err != nil { return Candidate{}, err }
	return Candidate{Kind: Unit, Value: result.Value, Format: result.Format, Confidence: result.Confidence, Ambiguous: result.Ambiguous, Assumptions: result.Assumptions, Candidates: result.Candidates}, nil
}

// DefaultRegistry returns a registry containing GoKit's built-in interpreters.
func DefaultRegistry() *Registry {
	registry := NewRegistry()
	registry.Register(expressionInterpreter{})
	registry.Register(unitInterpreter{})
	registry.Register(dateInterpreter{})
	registry.Register(numberInterpreter{})
	return registry
}
