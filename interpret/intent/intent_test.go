package intent

import (
	"math"
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/age"
	"github.com/codedbyhassan/gokit/calculator"
)

func TestParseCalculation(t *testing.T) {
	result, err := ParseAt("what's 20% of 500?", time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	calculated, ok := result.Value.(calculator.Result)
	if !ok {
		t.Fatalf("expected calculator result, got %T", result.Value)
	}
	if calculated.Value != 100 {
		t.Fatalf("expected 100, got %v", calculated.Value)
	}
}

func TestParseConversion(t *testing.T) {
	result, err := Parse("convert 10 miles to km")
	if err != nil {
		t.Fatal(err)
	}
	if result.Kind != Convert {
		t.Fatalf("expected convert intent, got %q", result.Kind)
	}
}

func TestParseAge(t *testing.T) {
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	result, err := ParseAt("how old is someone born 11-11-2011", asOf)
	if err != nil {
		t.Fatal(err)
	}
	if result.Kind != Age {
		t.Fatalf("expected age intent, got %q", result.Kind)
	}
	calculated, ok := result.Value.(age.Result)
	if !ok {
		t.Fatalf("expected age result, got %T", result.Value)
	}
	if calculated.Years != 14 || calculated.Months != 9 || calculated.Days != 22 {
		t.Fatalf("unexpected age: %+v", calculated)
	}
}

func TestParseQuantityCalculation(t *testing.T) {
	result, err := Parse("calculate 5kg plus 200g")
	if err != nil {
		t.Fatal(err)
	}
	calculated, ok := result.Value.(calculator.Result)
	if !ok {
		t.Fatalf("expected calculator result, got %T", result.Value)
	}
	if math.Abs(calculated.Value-5.2) > 1e-9 {
		t.Fatalf("expected 5.2, got %v", calculated.Value)
	}
}

func TestParseRejectsUnknownRequest(t *testing.T) {
	if _, err := Parse("tell me something interesting"); err == nil {
		t.Fatal("expected unknown request to fail")
	}
}
