package intent

import (
	"math"
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/age"
	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret/unit"
)

func TestParseCalculationWrappers(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"calculate 10 plus 5", 15},
		{"what is 10 multiplied by 5", 50},
		{"20% of 500", 100},
		{"half of 100", 50},
	}

	for _, tt := range tests {
		result, err := ParseAt(tt.input, time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatalf("%q: %v", tt.input, err)
		}
		calculated, ok := result.Value.(calculator.Result)
		if !ok || math.Abs(calculated.Value-tt.want) > 1e-9 {
			t.Fatalf("%q: expected %v, got %#v", tt.input, tt.want, result.Value)
		}
	}
}

func TestParseConversionValue(t *testing.T) {
	result, err := Parse("convert 10 miles to km")
	if err != nil {
		t.Fatal(err)
	}
	conversion, ok := result.Value.(unit.Conversion)
	if !ok {
		t.Fatalf("expected conversion, got %T", result.Value)
	}
	if math.Abs(conversion.Value-16.09344) > 1e-9 {
		t.Fatalf("expected 16.09344, got %v", conversion.Value)
	}
}

func TestParseAgeDateFormats(t *testing.T) {
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	for _, input := range []string{
		"how old is someone born 2011-11-11",
		"how old is someone born 11 November 2011",
	} {
		result, err := ParseAt(input, asOf)
		if err != nil {
			t.Fatalf("%q: %v", input, err)
		}
		calculated := result.Value.(age.Result)
		if calculated.Years != 14 || calculated.Months != 9 || calculated.Days != 22 {
			t.Fatalf("%q: unexpected age %+v", input, calculated)
		}
	}
}
