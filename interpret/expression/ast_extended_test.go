package expression

import (
	"math"
	"testing"
)

func TestEvaluatePercentAST(t *testing.T) {
	tests:=[]struct{input string;want float64}{
		{"25% of 800",200},
		{"20% increase on 500",600},
		{"20 percent increase on 500",600},
	}
	for _,tt:=range tests{n,err:=ParsePercentAST(tt.input);if err!=nil{t.Fatalf("%q: %v",tt.input,err)};got,err:=Evaluate(n);if err!=nil{t.Fatalf("%q: %v",tt.input,err)};if math.Abs(got-tt.want)>1e-9{t.Fatalf("%q: got %v want %v",tt.input,got,tt.want)}}
}

func TestEvaluateASTDivideByZero(t *testing.T){n,err:=ParseAST("10 / 0");if err!=nil{t.Fatal(err)};if _,err=Evaluate(n);err==nil{t.Fatal("expected division by zero")}}

func TestEvaluateASTNestedPrecedence(t *testing.T){n,err:=ParseAST("2 + 3 * (4 - 1)");if err!=nil{t.Fatal(err)};got,err:=Evaluate(n);if err!=nil{t.Fatal(err)};if got!=11{t.Fatalf("got %v want 11",got)}}
