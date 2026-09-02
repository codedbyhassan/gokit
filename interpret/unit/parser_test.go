package unit

import "testing"

func TestParseQuantity(t *testing.T) {
	tests := []struct { input string; value float64; symbol string }{
		{"5 kg", 5, "kg"},
		{"5kg", 5, "kg"},
		{"10 kilometers", 10, "km"},
		{"-2.5 lb", -2.5, "lb"},
		{"25 celsius", 25, "°C"},
		{"25°C", 25, "°C"},
		{"100 miles per hour", 100, "mph"},
	}
	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }
		q, ok := result.Value.(Quantity)
		if !ok { t.Fatalf("Parse(%q) returned %T, want Quantity", tt.input, result.Value) }
		if q.Value != tt.value || q.Unit.Symbol != tt.symbol { t.Fatalf("Parse(%q) = %#v, want %v %s", tt.input, q, tt.value, tt.symbol) }
	}
}

func TestParseConversion(t *testing.T) {
	tests := []struct { input string; want float64 }{
		{"5kg to pounds", 11.02311310924388},
		{"20 km to miles", 12.42742384474606},
		{"25 celsius to fahrenheit", 77},
		{"25°C in °F", 77},
	}
	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }
		c, ok := result.Value.(Conversion)
		if !ok { t.Fatalf("Parse(%q) returned %T, want Conversion", tt.input, result.Value) }
		if abs(c.Value-tt.want) > 1e-9 { t.Fatalf("Parse(%q) = %v, want %v", tt.input, c.Value, tt.want) }
	}
}

func TestParseRejectsIncompatibleUnits(t *testing.T) {
	if _, err := Parse("5 kg to miles"); err == nil { t.Fatal("expected incompatible-unit error") }
}

func abs(v float64) float64 { if v < 0 { return -v }; return v }
