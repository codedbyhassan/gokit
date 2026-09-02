package expression

import (
	"math"
	"testing"
)

func TestQuantityExpressions(t *testing.T) {
	tests := []struct { input string; want float64; symbol string }{
		{"5kg + 200g", 5.2, "kg"},
		{"5kg + 200g + 300mg", 5.2003, "kg"},
		{"10 miles + 5 km", 13.106855, "mi"},
		{"2m * 4", 8, "m"},
		{"5kg / 2", 2.5, "kg"},
		{"(2m + 3m) * 2", 10, "m"},
	}
	for _, tt := range tests {
		n, err := ParseQuantityAST(tt.input); if err != nil { t.Fatalf("%q: %v", tt.input, err) }
		got, err := EvaluateQuantity(n); if err != nil { t.Fatalf("%q: %v", tt.input, err) }
		if math.Abs(got.Value-tt.want) > 1e-6 { t.Errorf("%q: got %v want %v", tt.input, got.Value, tt.want) }
		if got.Unit.Symbol != tt.symbol { t.Errorf("%q: unit %q want %q", tt.input, got.Unit.Symbol, tt.symbol) }
	}
}

func TestQuantityParserKeepsScalars(t *testing.T) {
	n, err := ParseQuantityAST("2 * 4 + 3"); if err != nil { t.Fatal(err) }
	got, err := EvaluateQuantity(n); if err != nil { t.Fatal(err) }
	if got.Value != 11 { t.Fatalf("got %v want 11", got.Value) }
}

func TestQuantityDimensions(t *testing.T) {
	n, err := ParseQuantityAST("5kg + 2m"); if err != nil { t.Fatal(err) }
	if _, err = EvaluateQuantity(n); err == nil { t.Fatal("expected incompatible dimensions error") }
}

func TestQuantityDivisionByZero(t *testing.T) {
	n, err := ParseQuantityAST("5kg / 0"); if err != nil { t.Fatal(err) }
	if _, err = EvaluateQuantity(n); err == nil { t.Fatal("expected division by zero error") }
}
