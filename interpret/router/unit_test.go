package router

import (
	"testing"

	"github.com/codedbyhassan/gokit/interpret/unit"
)

func TestRoutesUnits(t *testing.T) {
	result, err := Parse("20 km to miles")
	if err != nil { t.Fatalf("Parse: %v", err) }
	if result.Kind != Unit { t.Fatalf("Kind = %q, want %q", result.Kind, Unit) }

	conversion, ok := result.AsUnitConversion()
	if !ok { t.Fatal("expected unit conversion") }
	if conversion.Value < 12.42 || conversion.Value > 12.43 {
		t.Fatalf("conversion = %v, want about 12.4274", conversion.Value)
	}
}

func TestRoutesQuantity(t *testing.T) {
	result, err := Parse("5 kg")
	if err != nil { t.Fatalf("Parse: %v", err) }
	quantity, ok := result.AsUnitQuantity()
	if !ok { t.Fatal("expected unit quantity") }
	if quantity.Value != 5 || quantity.Unit.Symbol != "kg" { t.Fatalf("unexpected quantity: %#v", quantity) }
}

var _ unit.Quantity
