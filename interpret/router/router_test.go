package router

import (
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/unit"
)

func TestParseRoutesByDomain(t *testing.T) {
	tests := []struct { name, input string; kind Kind }{
		{"number", "1,000", Number},
		{"date", "11-11-2011", Date},
		{"expression", "what is 10 multiplied by 5", Expression},
		{"percent", "25% of 800", Expression},
		{"word expression", "ten plus five", Expression},
		{"quantity", "5 kg", Unit},
		{"conversion", "20 km to miles", Unit},
		{"natural conversion", "what is 5 kg in pounds", Unit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil { t.Fatalf("Parse(%q) error = %v", tt.input, err) }
			if result.Kind != tt.kind { t.Fatalf("Parse(%q) kind = %q, want %q", tt.input, result.Kind, tt.kind) }
		})
	}
}

func TestParseTypedAccessors(t *testing.T) {
	numberResult, err := Parse("2.5 million")
	if err != nil { t.Fatal(err) }
	if value, ok := numberResult.AsNumber(); !ok || value != 2_500_000 { t.Fatalf("unexpected number result: %v, %v", value, ok) }

	dateResult, err := Parse("2011-11-11")
	if err != nil { t.Fatal(err) }
	value, ok := dateResult.AsDate()
	if !ok || !value.Equal(time.Date(2011, 11, 11, 0, 0, 0, 0, time.UTC)) { t.Fatalf("unexpected date result: %v, %v", value, ok) }

	expressionResult, err := Parse("10 plus 5")
	if err != nil { t.Fatal(err) }
	op, ok := expressionResult.AsExpression()
	if !ok || op.Left != 10 || op.Right != 5 { t.Fatalf("unexpected expression result: %+v, %v", op, ok) }

	unitResult, err := Parse("5 kg to pounds")
	if err != nil { t.Fatal(err) }
	conversion, ok := unitResult.Value.(unit.Conversion)
	if !ok || conversion.To.Symbol != "lb" { t.Fatalf("unexpected unit result: %#v, %v", conversion, ok) }
}

func TestParseUnknown(t *testing.T) {
	_, err := Parse("this is not supported")
	if err == nil { t.Fatal("expected an error") }
	if !interpretErrIs(err, interpret.ErrUnrecognizedInput) { t.Fatalf("error = %v, want ErrUnrecognizedInput", err) }
}

func interpretErrIs(err, target error) bool {
	for err != nil {
		if err == target { return true }
		unwrapper, ok := err.(interface{ Unwrap() error })
		if !ok { return false }
		err = unwrapper.Unwrap()
	}
	return false
}
