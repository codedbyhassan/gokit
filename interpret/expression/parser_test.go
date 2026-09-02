package expression

import "testing"

func TestParseExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  float64
		right float64
		op    Operator
	}{
		{"10 + 5", 10, 5, Add},
		{"10 plus 5", 10, 5, Add},
		{"10 added to 5", 10, 5, Add},
		{"10 - 3", 10, 3, Subtract},
		{"10 minus 3", 10, 3, Subtract},
		{"10 multiplied by 5", 10, 5, Multiply},
		{"10 times 5", 10, 5, Multiply},
		{"100 divided by 4", 100, 4, Divide},
		{"what is 10 multiplied by 5", 10, 5, Multiply},
		{"25% of 800", 800, 25, PercentOf},
		{"half of 100", 100, 2, Divide},
		{"20% increase on 500", 500, 20, PercentIncrease},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", tt.input, err)
		}
		if result.Value.Left != tt.left || result.Value.Right != tt.right || result.Value.Operator != tt.op {
			t.Errorf("Parse(%q) = %+v, want left=%v right=%v op=%q", tt.input, result.Value, tt.left, tt.right, tt.op)
		}
	}
}

func TestCalculateExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"what is 10 multiplied by 5", 50},
		{"100 divided by 4", 25},
		{"25% of 800", 200},
		{"half of 100", 50},
		{"20% increase on 500", 600},
	}

	for _, tt := range tests {
		result, err := Calculate(tt.input)
		if err != nil {
			t.Fatalf("Calculate(%q) error = %v", tt.input, err)
		}
		if result.Value != tt.want {
			t.Errorf("Calculate(%q) = %v, want %v", tt.input, result.Value, tt.want)
		}
	}
}

func TestParseInvalidExpression(t *testing.T) {
	if _, err := Parse("do something weird"); err == nil {
		t.Fatal("expected an error")
	}
}

func TestParsePercentIncreaseVariants(t *testing.T) {
	tests := []string{"20% increase on 500", "20 percent increase on 500"}
	for _, input := range tests {
		result, err := Parse(input)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", input, err)
		}
		if result.Value.Operator != PercentIncrease || result.Value.Left != 500 || result.Value.Right != 20 {
			t.Fatalf("Parse(%q) = %+v", input, result.Value)
		}
	}
}
