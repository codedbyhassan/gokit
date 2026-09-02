package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/pipeline"
	"github.com/codedbyhassan/gokit/interpret/router"
)

func TestPipelineCrossPackageFlows(t *testing.T) {
	tests := []struct {
		name string
		input string
		wantKind string
		wantValue float64
	}{
		{"expression", "what is 20% of 500", "calculate", 100},
		{"number", "1.5k", "number", 1500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pipeline.Parse(tt.input)
			if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }
			if result.Kind != tt.wantKind { t.Fatalf("kind=%q want %q", result.Kind, tt.wantKind) }
			if tt.name == "expression" {
				got, ok := result.Value.(calculator.Result)
				if !ok { t.Fatalf("value type=%T want calculator.Result", result.Value) }
				if got.Value != tt.wantValue { t.Fatalf("value=%v want %v", got.Value, tt.wantValue) }
			} else if result.Value != tt.wantValue { t.Fatalf("value=%v want %v", result.Value, tt.wantValue) }
			if len(result.Plan.Steps) == 0 { t.Fatal("expected execution plan") }
		})
	}
}

func TestPipelineDateAndAmbiguity(t *testing.T) {
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	result, err := pipeline.ParseAt("how old is someone born 11-11-2011", asOf)
	if err != nil { t.Fatalf("date intent: %v", err) }
	if result.Kind != "age" || result.Plan.Intent != "age" { t.Fatalf("unexpected result: %+v", result) }

	_, err = pipeline.Parse("05-06-2020")
	if !errors.Is(err, interpret.ErrAmbiguousInput) { t.Fatalf("expected ambiguity, got %v", err) }
}

func TestRegistryIgnoresNilInterpreter(t *testing.T) {
	registry := router.NewRegistry()
	registry.Register(nil)
	registry.Register(nil)
	if _, err := registry.Parse("1.5k"); err == nil { t.Fatal("empty registry should not parse input") }
}

type testInterpreter struct{ candidate router.Candidate }
func (i testInterpreter) Interpret(string) (router.Candidate, error) { return i.candidate, nil }

func TestRegistryConfidenceAndRegistrationOrder(t *testing.T) {
	registry := router.NewRegistry()
	registry.Register(testInterpreter{router.Candidate{Kind: router.Number, Value: "first", Confidence: 0.8}})
	registry.Register(testInterpreter{router.Candidate{Kind: router.Date, Value: "second", Confidence: 0.9}})
	result, err := registry.Parse("anything")
	if err != nil { t.Fatalf("Parse: %v", err) }
	if result.Kind != router.Date || result.Value != "second" { t.Fatalf("unexpected winner: %+v", result) }
}
