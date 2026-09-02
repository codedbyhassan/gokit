package pipeline

import (
	"errors"
	"testing"
	"time"

	"github.com/codedbyhassan/gokit/calculator"
	"github.com/codedbyhassan/gokit/interpret"
)

func TestParseAtUsesIntentAndExecutesNaturalCommands(t *testing.T) {
	asOf := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	result, err := ParseAt("what is 20% of 500", asOf)
	if err != nil { t.Fatal(err) }
	if result.Source != IntentSource { t.Fatalf("source = %q, want %q", result.Source, IntentSource) }
	if result.Kind != "calculate" { t.Fatalf("kind = %q, want calculate", result.Kind) }
	value, ok := result.Value.(calculator.Result)
	if !ok { t.Fatalf("value type = %T, want calculator.Result", result.Value) }
	if value.Value != 100 { t.Fatalf("value = %v, want 100", value.Value) }
	if len(result.Plan.Steps) != 2 { t.Fatalf("plan steps = %d, want 2", len(result.Plan.Steps)) }
	if result.Interpretation == nil { t.Fatal("interpretation should be preserved") }
}

func TestParseAtFallsBackToRouter(t *testing.T) {
	result, err := ParseAt("1.5k", time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
	if err != nil { t.Fatal(err) }
	if result.Source != RouterSource { t.Fatalf("source = %q, want %q", result.Source, RouterSource) }
	if result.Kind != "number" { t.Fatalf("kind = %q, want number", result.Kind) }
	if result.Value.(float64) != 1500 { t.Fatalf("value = %v, want 1500", result.Value) }
}

func TestParseAtPreservesContextualErrors(t *testing.T) {
	_, err := ParseAt("how old is someone born 31-02-2020", time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
	if err == nil { t.Fatal("expected invalid birth date error") }
	if errors.Is(err, interpret.ErrUnrecognizedInput) { t.Fatal("invalid age input should not fall through as unrecognized") }
}

func TestParseAtRejectsEmptyInput(t *testing.T) {
	_, err := ParseAt("   ", time.Now())
	if !errors.Is(err, interpret.ErrEmptyInput) { t.Fatalf("err = %v, want ErrEmptyInput", err) }
}
