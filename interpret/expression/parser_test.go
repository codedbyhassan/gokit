package expression

import (
	"math"
	"testing"
)

func TestParseExpressions(t *testing.T) {
	tests := []struct { input string; left, right float64; op Operator }{
		{"10 + 5", 10, 5, Add}, {"10 plus 5", 10, 5, Add}, {"10 added to 5", 10, 5, Add},
		{"10 - 3", 10, 3, Subtract}, {"10 minus 3", 10, 3, Subtract}, {"10 multiplied by 5", 10, 5, Multiply},
		{"10 times 5", 10, 5, Multiply}, {"100 divided by 4", 100, 4, Divide}, {"what is 10 multiplied by 5", 10, 5, Multiply},
		{"25% of 800", 800, 25, PercentOf}, {"half of 100", 100, 2, Divide}, {"20% increase on 500", 500, 20, PercentIncrease},
	}
	for _, tt := range tests { result, err := Parse(tt.input); if err != nil { t.Fatalf("Parse(%q): %v", tt.input, err) }; if result.Value.Left != tt.left || result.Value.Right != tt.right || result.Value.Operator != tt.op { t.Errorf("Parse(%q) = %+v", tt.input, result.Value) } }
}

func TestCalculateExpressions(t *testing.T) {
	tests := []struct { input string; want float64 }{{"what is 10 multiplied by 5",50},{"100 divided by 4",25},{"25% of 800",200},{"half of 100",50},{"20% increase on 500",600},{"-10 + 4",-6},{"10 subtract from 5",-5}}
	for _, tt := range tests { result, err := Calculate(tt.input); if err != nil { t.Fatalf("Calculate(%q): %v", tt.input, err) }; if math.Abs(result.Value-tt.want) > 1e-9 { t.Errorf("Calculate(%q) = %v, want %v", tt.input, result.Value, tt.want) } }
}

func TestParseInvalidExpression(t *testing.T) { if _, err := Parse("do something weird"); err == nil { t.Fatal("expected an error") } }

func TestParsePercentIncreaseVariants(t *testing.T) {
	for _, input := range []string{"20% increase on 500", "20 percent increase on 500"} { result, err := Parse(input); if err != nil { t.Fatalf("Parse(%q): %v", input, err) }; if result.Value.Operator != PercentIncrease || result.Value.Left != 500 || result.Value.Right != 20 { t.Fatalf("Parse(%q) = %+v", input, result.Value) } }
}
