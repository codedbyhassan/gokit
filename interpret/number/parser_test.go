package number

import "testing"

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"1000", 1000},
		{"1,000", 1000},
		{"1 000", 1000},
		{"1.5", 1.5},
		{"1.5k", 1500},
		{"5k", 5000},
		{"2 million", 2000000},
		{"2.5 million", 2500000},
		{"five", 5},
		{"twenty five", 25},
		{"one hundred", 100},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", tt.input, err)
		}
		if result.Value != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Value, tt.want)
		}
		if result.Confidence <= 0 {
			t.Errorf("Parse(%q) returned non-positive confidence", tt.input)
		}
	}
}

func TestParseInvalidNumber(t *testing.T) {
	if _, err := Parse("not a number"); err == nil {
		t.Fatal("expected an error")
	}
}
