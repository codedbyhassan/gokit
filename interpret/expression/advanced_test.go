package expression

import "testing"

func TestParseNaturalOperands(t *testing.T) {
	tests := []struct {
		input string
		left  float64
		right float64
		op    Operator
	}{
		{"ten plus five", 10, 5, Add},
		{"twenty five times four", 25, 4, Multiply},
		{"one hundred divided by five", 100, 5, Divide},
		{"2.5 million plus 500 thousand", 2500000, 500000, Add},
		{"-10 plus 4", -10, 4, Add},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", tt.input, err)
		}
		if result.Value.Left != tt.left || result.Value.Right != tt.right || result.Value.Operator != tt.op {
			t.Errorf("Parse(%q) = %+v, want %v %v %q", tt.input, result.Value, tt.left, tt.right, tt.op)
		}
	}
}

func TestParsePunctuation(t *testing.T) {
	result, err := Parse("What is 1,000 multiplied by 2?")
	if err != nil {
		t.Fatal(err)
	}
	if result.Value.Left != 1000 || result.Value.Right != 2 || result.Value.Operator != Multiply {
		t.Fatalf("unexpected result: %+v", result.Value)
	}
}
