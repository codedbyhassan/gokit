package unit

import "testing"

func TestParseQuantity(t *testing.T) {
	tests := []struct {
		input string
		value float64
		symbol string
	}{
		{"5 kg", 5, "kg"},
		{"10 kilometers", 10, "km"},
		{"25 celsius", 25, "°C"},
		{"100 miles per hour", 100, "mph"},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }
		quantity, ok := result.Value.(Quantity)
		if !ok { t.Fatalf("Parse(%q): expected Quantity", tt.input) }
		if quantity.Value != tt.value || quantity.Unit.Symbol != tt.symbol {
			t.Fatalf("Parse(%q) = %#v, want %v %s", tt.input, quantity, tt.value, tt.symbol)
		}
	}
}

func TestParseConversion(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"5kg to pounds", 11.02311310924388},
		{"20 km to miles", 12.42742384474606},
		{"25 celsius to fahrenheit", 77},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }
		conversion, ok := result.Value.(Conversion)
		if !ok { t.Fatalf("Parse(%q): expected Conversion", tt.input) }
		if conversion.Value != tt.want {
			t.Fatalf("Parse(%q) = %v, want %v", tt.input, conversion.Value, tt.want)
		}
	}
}

func TestParseRejectsIncompatibleUnits(t *testing.T) {
	if _, err := Parse("5 kg to miles"); err == nil {
		t.Fatal("expected incompatible-unit error")
	}
}
